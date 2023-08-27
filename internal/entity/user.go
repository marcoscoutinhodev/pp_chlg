package entity

type User struct {
	ID                      string `bson:"_id,omitempty"`
	UserID                  string `bson:"user_id"`
	FirstName               string `bson:"first_name"`
	LastName                string `bson:"last_name"`
	Email                   string `bson:"email"`
	Password                string `bson:"password"`
	TaxpayeerIdentification string `bson:"taxpayeer_identification"`
	Group                   string `bson:"group"`
}

func NewUser(firstName, lastName, email, password, taxpayeerIdentification, group string) *User {
	return &User{
		FirstName:               firstName,
		LastName:                lastName,
		Email:                   email,
		Password:                password,
		TaxpayeerIdentification: taxpayeerIdentification,
		Group:                   group,
	}
}
