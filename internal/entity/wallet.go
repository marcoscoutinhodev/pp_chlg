package entity

type Wallet struct {
	ID      string  `bson:"_id,omitempty"`
	UserID  string  `bson:"user_id"`
	Balance float64 `bson:"balance"`
}

func NewWallet(userID string, balance float64) *Wallet {
	return &Wallet{
		UserID:  userID,
		Balance: balance,
	}
}
