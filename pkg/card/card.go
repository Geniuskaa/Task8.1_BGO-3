package card

import (
	"github.com/Geniuskaa/Task8.1_BGO-3/pkg/transaction"
	"time"
)

type Card struct {
	Id int64
	Issuer string
	Currency string
	Balance int64
	Number string
	Transactions []*transaction.Transaction
}

type Service struct {
	bank string
	StoreOfCards []*Card
}

type part struct {
	monthTimestamp int64
	transactions []*transaction.Transaction
}

func NewService(bankName string) *Service {
	return &Service{
		bank: bankName,
		StoreOfCards: []*Card{}}
}

func (s *Service) AddCard(id int64, issuer string, currency string, balance int64, number string) {
	s.StoreOfCards = append(s.StoreOfCards, &Card{
		Id:       id,
		Issuer:   issuer,
		Currency: currency,
		Balance:  balance,
		Number:   number,
	})
}

func (s *Card) AddTransaction(amount int64, MCC string, time time.Time) {
	s.Transactions = append(s.Transactions, &transaction.Transaction{
		Id:     20,
		Amount: amount * 100,
		MCC:    MCC,
		Date:   time,
		Status: "Completed",
	})
}









