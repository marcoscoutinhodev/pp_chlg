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

var createUserInputMock = &InputUserAuthenticationCreateUserDTO{
	FirstName:               "any_first_name",
	LastName:                "any_last_name",
	Email:                   "any_email",
	Password:                "any_password",
	TaxpayeerIdentification: "any_taxpayeer_identification",
	Group:                   "any_group",
}

func (s *UserAuthenticationCreateUserSuite) TestGivenInvalidInput_ShouldReturnValidationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayeerIdentification,
		createUserInputMock.Group,
	)

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("Validate", *userEntityMock).Return([]string{"any_error", "other_error"})

	sut := NewUserAuthentication(&mocks.IdentityManagerMock{}, userValidatorMock, &mocks.UserRepositoryMock{})

	output, err := sut.CreateUser(context.Background(), createUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationCreateUserDTO{
		ValidationError:  []string{"any_error", "other_error"},
		DuplicationError: "",
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
		createUserInputMock.Group,
	)

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("Validate", *userEntityMock).Return(nil)

	userRepositoryMock := &mocks.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserIsRegistered", context.Background(), *userEntityMock).Return(true, nil)

	sut := NewUserAuthentication(&mocks.IdentityManagerMock{}, userValidatorMock, userRepositoryMock)

	output, err := sut.CreateUser(context.Background(), createUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationCreateUserDTO{
		ValidationError:  nil,
		DuplicationError: "email or and taxpayeer identification is already registered",
	}, *output)

	userValidatorMock.AssertExpectations(s.T())
	userRepositoryMock.AssertExpectations(s.T())
}

func (s *UserAuthenticationCreateUserSuite) TestGivenInvalidGroup_ShouldReturnValidationError() {
	userEntityMock := entity.NewUser(
		createUserInputMock.FirstName,
		createUserInputMock.LastName,
		createUserInputMock.Email,
		createUserInputMock.Password,
		createUserInputMock.TaxpayeerIdentification,
		createUserInputMock.Group,
	)

	identityManagerMock := &mocks.IdentityManagerMock{}
	identityManagerMock.On("GetGroupID", context.Background(), userEntityMock.Group).Return("", nil)

	userValidatorMock := &mocks.UserValidatorMock{}
	userValidatorMock.On("Validate", *userEntityMock).Return(nil)

	userRepositoryMock := &mocks.UserRepositoryMock{}
	userRepositoryMock.On("CheckUserIsRegistered", context.Background(), *userEntityMock).Return(false, nil)

	sut := NewUserAuthentication(identityManagerMock, userValidatorMock, userRepositoryMock)

	output, err := sut.CreateUser(context.Background(), createUserInputMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OutputUserAuthenticationCreateUserDTO{
		ValidationError:  []string{"invalid group provided"},
		DuplicationError: "",
	}, *output)

	identityManagerMock.AssertExpectations(s.T())
	userValidatorMock.AssertExpectations(s.T())
	userRepositoryMock.AssertExpectations(s.T())
}

func TestUserAuthentication_CreateUserSuite(t *testing.T) {
	suite.Run(t, new(UserAuthenticationCreateUserSuite))
}
