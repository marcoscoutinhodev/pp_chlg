package entity

type User struct {
	ID                     string `bson:"_id,omitempty"`
	UserID                 string `bson:"user_id"`
	FirstName              string `bson:"first_name"`
	LastName               string `bson:"last_name"`
	Email                  string `bson:"email"`
	Password               string `bson:"-"`
	TaxpayerIdentification string `bson:"taxpayer_identification"`
	Role                   string `bson:"role"`
}

func NewUser(firstName, lastName, email, password, taxpayerIdentification, role string) *User {
	return &User{
		FirstName:              firstName,
		LastName:               lastName,
		Email:                  email,
		Password:               password,
		TaxpayerIdentification: taxpayerIdentification,
		Role:                   role,
	}
}
