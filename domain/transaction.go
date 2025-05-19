package domain

import "github.com/google/uuid"

type Transaction struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
	From   string `json:"from"`
	To     string `json:"to"`
}

func NewTransaction(amount int64, from string, to string) *Transaction {
	t := new(Transaction)
	t.ID = uuid.NewString()
	t.Amount = amount
	t.From = from
	t.To = to
	return t
}
