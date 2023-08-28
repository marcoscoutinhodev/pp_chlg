package usecase

import (
	"context"
	"testing"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/marcoscoutinhodev/pp_chlg/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserAuthenticationCreateUserSuite struct {
	suite.Suite
}

var createUserInputMock = &UserAuthentication_CreateUserInputDTO{
	FirstName:               "any_first_name",
	LastName:                "any_last_name",
	Email:                   "any_email",
	Password:                "any_password",
	TaxpayeerIdentification: "any_taxpayeer_identification",
	Role:                    "any_role",
}

func (s *UserAuthenticationCreateUserSuite) TestGivenInvalidInput_ShouldReturnValidationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayeerIdentification,
		createUserInputMock.Role,
	)

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("Validate", *userEntityMock).Return([]string{"any_error", "other_error"})

	sut := NewUserAuthentication(&mocks.IdentityManagerMock{}, userValidatorMock, &mocks.UserRepositoryMock{})

	output, err := sut.CreateUser(context.Background(), createUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationDTO{
		StatusCode: 400,
		Success:    false,
		Errors:     []string{"any_error", "other_error"},
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
}

func (s *UserAuthenticationCreateUserSuite) TestGivenEmailOrTaxpayeerIdentificationRegistered_ShouldReturnDuplicationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayeerIdentification,
		createUserInputMock.Role,
	)

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("Validate", *userEntityMock).Return(nil)

	userRepositoryMock := &mocks.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserIsRegistered", context.Background(), *userEntityMock).Return(true, nil)

	sut := NewUserAuthentication(&mocks.IdentityManagerMock{}, userValidatorMock, userRepositoryMock)

	output, err := sut.CreateUser(context.Background(), createUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationDTO{
		StatusCode: 400,
		Success:    false,
		Errors:     []string{"email and/or taxpayeer identification are already registered"},
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
	userRepositoryMock.AssertExpectations(s.T())
}

func TestUserAuthentication_CreateUserSuite(t *testing.T) {
	suite.Run(t, new(UserAuthenticationCreateUserSuite))
}
