package entity

import (
	"strings"
	"time"
)

type Transfer struct {
	ID     string  `bson:"_id,omitempty"`
	Payer  string  `bson:"payer"`
	Payee  string  `bson:"payee"`
	Amount float64 `bson:"amount"`
	Date   string  `bson:"date"`
}

func NewTransfer(payer, payee string, amount float64) *Transfer {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	now := time.Now().In(loc)
	nowAsString := strings.Replace(now.Format(time.RFC3339), "T", " ", 1)

	return &Transfer{
		Payer:  payer,
		Payee:  payee,
		Amount: amount,
		Date:   nowAsString,
	}
}
