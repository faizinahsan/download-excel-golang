package domain

import "time"

type Transaction struct {
	ID              int32
	Amount          float64
	Date            time.Time
	TransactionName string
	ReferenceNumber string
}
