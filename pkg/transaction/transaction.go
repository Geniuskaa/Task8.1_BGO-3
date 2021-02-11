package transaction

import "time"

type Transaction struct {
	Id int64 `json:"id"`
	Amount int64 `json:"amount"`
	MCC string `json:"mcc"`
	Date time.Time `json:"date"`
	Status string `json:"status"`
}




