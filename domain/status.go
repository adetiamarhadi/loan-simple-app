package domain

type statusType string

// enum for status loan application
const (
	Open     = "open"
	Approved = "approved"
	Rejected = "rejected"
)