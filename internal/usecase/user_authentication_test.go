package usecase

import (
	"context"
	"errors"
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
	FirstName:              "any_first_name",
	LastName:               "any_last_name",
	Email:                  "any_email",
	Password:               "any_password",
	TaxpayerIdentification: "any_taxpayer_identification",
	Role:                   "any_role",
}

func (s *UserAuthenticationCreateUserSuite) TestGivenInvalidInput_ShouldReturnValidationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayerIdentification,
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

func (s *UserAuthenticationCreateUserSuite) TestGivenEmailOrTaxpayerIdentificationRegistered_ShouldReturnDuplicationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayerIdentification,
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
		Errors:     []string{"email and/or taxpayer identification are already registered"},
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
	userRepositoryMock.AssertExpectations(s.T())
}

func TestUserAuthentication_CreateUserSuite(t *testing.T) {
	suite.Run(t, new(UserAuthenticationCreateUserSuite))
}

type UserAuthenticationAuthenticateUserSuite struct {
	suite.Suite
}

var authenticateUserInputMock = &UserAuthentication_AuthenticateUserInputDTO{
	Email:    "any_email",
	Password: "any_password",
}

func (s *UserAuthenticationAuthenticateUserSuite) TestGivenInvalidInput_ShouldReturnError() {

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("ValidateEmailAndPassword", authenticateUserInputMock.Email, authenticateUserInputMock.Password).Return([]string{"any_error", "other_error"})

	sut := NewUserAuthentication(&mocks.IdentityManagerMock{}, userValidatorMock, &mocks.UserRepositoryMock{})

	output, err := sut.AuthenticateUser(context.Background(), authenticateUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationDTO{
		StatusCode: 401,
		Success:    false,
		Errors:     []string{"email and/or password are invalid"},
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
}

func (s *UserAuthenticationAuthenticateUserSuite) TestGivenInvalidEmailOrPassword_ShouldReturnError() {
	identityManagerMock := &mocks.IdentityManagerMock{}
	identityManagerMock.On("AuthenticateUser", context.Background(), authenticateUserInputMock.Email, authenticateUserInputMock.Password).Return(nil, errors.New("any_error"))

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("ValidateEmailAndPassword", authenticateUserInputMock.Email, authenticateUserInputMock.Password).Return([]string{})

	sut := NewUserAuthentication(identityManagerMock, userValidatorMock, &mocks.UserRepositoryMock{})

	output, err := sut.AuthenticateUser(context.Background(), authenticateUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationDTO{
		StatusCode: 401,
		Success:    false,
		Errors:     []string{"email and/or password are invalid"},
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
}

func TestUserAuthentication_AuthenticateUserSuite(t *testing.T) {
	suite.Run(t, new(UserAuthenticationAuthenticateUserSuite))
}
