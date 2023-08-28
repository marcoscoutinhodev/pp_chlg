package entity

type User struct {
	ID                      string `bson:"_id,omitempty"`
	UserID                  string `bson:"user_id"`
	FirstName               string `bson:"first_name"`
	LastName                string `bson:"last_name"`
	Email                   string `bson:"email"`
	Password                string `bson:"-"`
	TaxpayeerIdentification string `bson:"taxpayeer_identification"`
	Role                    string `bson:"role"`
}

func NewUser(firstName, lastName, email, password, taxpayeerIdentification, role string) *User {
	return &User{
		FirstName:               firstName,
		LastName:                lastName,
		Email:                   email,
		Password:                password,
		TaxpayeerIdentification: taxpayeerIdentification,
		Role:                    role,
	}
}
