package mocks

import (
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type UserValidatorMock struct {
	mock.Mock
}

func (u *UserValidatorMock) Validate(user entity.User) (errors []string) {
	args := u.Called(user)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).([]string)
}

func (u *UserValidatorMock) ValidateEmailAndPassword(email, password string) (errors []string) {
	args := u.Called(email, password)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).([]string)
}
