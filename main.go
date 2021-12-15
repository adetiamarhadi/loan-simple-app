package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adetiamarhadi/loan-simple-app/domain"
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

func createLoan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "create loan")
}

func updateLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "update loan for application number : %s", ps.ByName("applicationNumber"))
}

func readLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loanApplication domain.LoanApplication

	loanApplication.FullName = "Adetia"
	loanApplication.BirthDate = time.Now()
	loanApplication.Gender = domain.Male
	loanApplication.MobileNumber = "6281200001111"
	loanApplication.Email = "adet@mail.com"
	loanApplication.ApplicationNumber = ps.ByName("applicationNumber")
	loanApplication.Status = domain.Open
	loanApplication.LoanAmount = 2000000
	loanApplication.LoanTerm = 12
	loanApplication.LoanInterest = 1.49
	loanApplication.LoanInterestMonthlyAmount = 0
	loanApplication.LoanAmountTotal = 0
	loanApplication.PaidMonthlyCount = 1
	loanApplication.NotPaidMonthlyCount = 12
	loanApplication.PaidAmount = 0
	loanApplication.NotPaidAmount = 2000000
	loanApplication.DueDatePayment = time.Now()
	loanApplication.LastPaymentDate = time.Now()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(loanApplication)
}

func approveLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "approve loan for application number : %s", ps.ByName("applicationNumber"))
}

func rejectLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "reject loan for application number : %s", ps.ByName("applicationNumber"))
}