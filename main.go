package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/adetiamarhadi/loan-simple-app/domain"
	"github.com/adetiamarhadi/loan-simple-app/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", homePage)
	router.POST("/loan", createLoan)
	router.GET("/loan/:applicationNumber", readLoan)
	router.POST("/loan/:applicationNumber", approveLoan)
	router.PUT("/loan/:applicationNumber", updateLoan)
	router.DELETE("/loan/:applicationNumber", rejectLoan)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "welcome")
}

func createLoan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// connect to db
	db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/core")
    if err != nil {
        log.Fatalln(err)
	}

	// mapping request to struct
	loanApplication := &domain.LoanApplication{}
	if err := populateModelFromHandler(w, r, params, loanApplication); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "err911: fail to process your request")
		return
	}

	loanApplication.Status = domain.Open

	totalInterest := loanApplication.LoanAmount * (loanApplication.LoanInterest / 100) * float64(loanApplication.LoanTerm)

	loanApplication.Total = loanApplication.LoanAmount + totalInterest

	loanApplication.MonthlyPayment = loanApplication.Total / float64(loanApplication.LoanTerm)

	loanApplication.ApplicationNumber = util.RandomAlphaNumeric(6)

	rows, err := db.NamedQuery("SELECT id, name, mobile_number, email FROM users WHERE email=:email OR mobile_number=:mobile_number", loanApplication)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "err914: fail to get data user")
		return
	}

	data := []domain.LoanApplication{}

	// add user to slice
	for rows.Next() {
		var la domain.LoanApplication
		err = rows.StructScan(&la)
		if err == nil {
			data = append(data, la)
		}
	}

	sizeOfData := len(data)

	if sizeOfData > 1 {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "err912: please change your mobile number or email address")
		return
	} else if sizeOfData == 1 {
		la := data[0]
		loanApplication.ID = la.ID

		var isPaidOff int

		err = db.Get(&isPaidOff, "SELECT count(id) FROM loans WHERE user_id = ? AND status = ?", loanApplication.ID, domain.PaidOff)

		if isPaidOff == 1 && la.Email == loanApplication.Email && la.MobileNumber == loanApplication.MobileNumber {
			tx := db.MustBegin()
			tx.NamedExec("INSERT INTO loans (user_id, code, amount, term, interest, monthly_payment, total, status) VALUES (:id, :code, :amount, :term, :interest, :monthly_payment, :total, :status)", loanApplication)
			tx.Commit()
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprint(w, "err913: your request still processing or your loan not yet paid off")
			return
		}
	} else {
		tx := db.MustBegin()
		tx.NamedExec("INSERT INTO users (name, mobile_number, email) VALUES (:name, :mobile_number, :email)", loanApplication)
		tx.NamedExec("INSERT INTO loans (user_id, code, amount, term, interest, monthly_payment, total, status) VALUES (LAST_INSERT_ID(), :code, :amount, :term, :interest, :monthly_payment, :total, :status)", loanApplication)
		tx.Commit()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(loanApplication)
}

func updateLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "update loan for application number : %s", ps.ByName("applicationNumber"))
}

func readLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// connect to db
	db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/core")
    if err != nil {
        log.Fatalln(err)
	}

	row := db.QueryRowx("SELECT l.id, u.name, u.mobile_number, u.email, l.code, l.status, l.amount, l.term, l.interest, l.monthly_payment, l.total FROM users u INNER JOIN loans l ON u.id = l.user_id WHERE l.code = ?", ps.ByName("applicationNumber"))

	var loan domain.LoanApplication

	err = row.StructScan(&loan)

	if err != nil && err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "err411: fail to get loan detail")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(loan)
}

func approveLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// connect to db
	db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/core")
    if err != nil {
        log.Fatalln(err)
	}

	row := db.QueryRowx("SELECT l.id, u.name, u.mobile_number, u.email, l.code, l.status, l.amount, l.term, l.interest, l.monthly_payment, l.total FROM users u INNER JOIN loans l ON u.id = l.user_id WHERE l.code = ?", ps.ByName("applicationNumber"))

	var loan domain.LoanApplication

	err = row.StructScan(&loan)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "err412: loan_id not found")
		return
	}

	tx := db.MustBegin()
	tx.NamedExec("UPDATE loans SET status = 'approved' WHERE code = :code", loan)
	tx.Commit()

	row = db.QueryRowx("SELECT l.id, u.name, u.mobile_number, u.email, l.code, l.status, l.amount, l.term, l.interest, l.monthly_payment, l.total FROM users u INNER JOIN loans l ON u.id = l.user_id WHERE l.code = ?", ps.ByName("applicationNumber"))

	err = row.StructScan(&loan)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(loan)
}

func rejectLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "reject loan for application number : %s", ps.ByName("applicationNumber"))
}

func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}