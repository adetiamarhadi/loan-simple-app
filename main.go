package main

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Fprintf(w, "read loan for application number : %s", ps.ByName("applicationNumber"))
}

func approveLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "approve loan for application number : %s", ps.ByName("applicationNumber"))
}

func rejectLoan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "reject loan for application number : %s", ps.ByName("applicationNumber"))
}