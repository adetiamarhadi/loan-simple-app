package domain

import "time"

// LoanApplication is data structure for loan application
type LoanApplication struct {
	FullName  string
	BirthDate time.Time
	Gender genderType
	MobileNumber string
	Email string
	ApplicationNumber string
	Status statusType
	LoanAmount float64
	LoanTerm int
	LoanInterest float64
	LoanInterestMonthlyAmount float64
	LoanAmountTotal float64
	PaidMonthlyCount int
	NotPaidMonthlyCount int
	PaidAmount float64
	NotPaidAmount float64
	DueDatePayment time.Time
	LastPaymentDate time.Time 
}