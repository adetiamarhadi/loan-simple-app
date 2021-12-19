package domain

// LoanApplication is data structure for loan application
type LoanApplication struct {
	ID int `db:"id" json:"-"`
	FullName  string `db:"name" json:"full_name"`
	MobileNumber string `db:"mobile_number" json:"mobile_number"`
	Email string `db:"email" json:"email"`
	ApplicationNumber string `db:"code" json:"loan_id"`
	Status statusType `db:"status" json:"status"`
	LoanAmount float64 `db:"amount" json:"loan_amount"`
	LoanTerm int `db:"term" json:"loan_term"`
	LoanInterest float64 `db:"interest" json:"loan_interest"`
	MonthlyPayment float64 `db:"monthly_payment" json:"monthly_payment"`
	Total float64 `db:"total" json:"total"`
}