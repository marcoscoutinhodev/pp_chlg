package entity

type Transfer struct {
	ID     string  `bson:"_id,omitempty"`
	Payer  string  `bson:"payer"`
	Payee  string  `bson:"payee"`
	Amount float64 `bson:"amount"`
}

func NewTransfer(payer, payee string, amount float64) *Transfer {
	return &Transfer{
		Payer:  payer,
		Payee:  payee,
		Amount: amount,
	}
}
