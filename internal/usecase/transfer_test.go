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

type TransferSuite struct {
	suite.Suite
}

var inputTransferMock = &InputTransferDTO{
	Payer:  "any_payer_id",
	Payee:  "any_payee_id",
	Amount: 10,
}

func (s *TransferSuite) TestGivenAmountGreaterThanBalance_ShouldReturnError() {
	walletRepositoryMock := &mocks.WalletRepositoryMock{}
	walletRepositoryMock.On("Load", context.Background(), "any_payer_id").Return(entity.NewWallet("any_payer_id", 0), nil)

	sut := NewTransfer(&mocks.TransferAuthorizationServiceMock{}, &mocks.EmailNotificationServiceMock{}, walletRepositoryMock)
	output, err := sut.Execute(context.Background(), inputTransferMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OuputTransferDTO{
		StatusCode: 402,
		Success:    false,
		Errors:     []string{"insufficient funds"},
	}, *output)

	walletRepositoryMock.AssertExpectations(s.T())
}

func (s *TransferSuite) TestGivenErrorInTransferAuthorizationService_ShouldReturnError() {
	walletRepositoryMock := &mocks.WalletRepositoryMock{}
	walletRepositoryMock.On("Load", context.Background(), "any_payer_id").Return(entity.NewWallet("any_payer_id", 10), nil)

	transferMock := entity.NewTransfer(inputTransferMock.Payer, inputTransferMock.Payee, inputTransferMock.Amount)
	transferAuthorizationServiceMock := &mocks.TransferAuthorizationServiceMock{}
	transferAuthorizationServiceMock.On("Check", context.Background(), *transferMock).Return(errors.New("unauthorized"))

	sut := NewTransfer(transferAuthorizationServiceMock, &mocks.EmailNotificationServiceMock{}, walletRepositoryMock)
	output, err := sut.Execute(context.Background(), inputTransferMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OuputTransferDTO{
		StatusCode: 422,
		Success:    false,
		Errors:     []string{"the transfer was not authorized"},
	}, *output)

	walletRepositoryMock.AssertExpectations(s.T())
	transferAuthorizationServiceMock.AssertExpectations(s.T())
}

func (s *TransferSuite) TestGivenErrorInWalletRepository_Transfer_ShouldReturnError() {
	transferMock := entity.NewTransfer(inputTransferMock.Payer, inputTransferMock.Payee, inputTransferMock.Amount)

	walletRepositoryMock := &mocks.WalletRepositoryMock{}
	walletRepositoryMock.On("Load", context.Background(), "any_payer_id").Return(entity.NewWallet("any_payer_id", 10), nil)
	walletRepositoryMock.On("Transfer", context.Background(), *transferMock).Return(nil, nil, errors.New("any_error"))

	transferAuthorizationServiceMock := &mocks.TransferAuthorizationServiceMock{}
	transferAuthorizationServiceMock.On("Check", context.Background(), *transferMock).Return(nil)

	sut := NewTransfer(transferAuthorizationServiceMock, &mocks.EmailNotificationServiceMock{}, walletRepositoryMock)
	output, err := sut.Execute(context.Background(), inputTransferMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OuputTransferDTO{
		StatusCode: 422,
		Success:    false,
		Errors:     []string{"unexpected error completing transfer"},
	}, *output)

	walletRepositoryMock.AssertExpectations(s.T())
	transferAuthorizationServiceMock.AssertExpectations(s.T())
}

func (s *TransferSuite) TestGivenErrorInEmailNotification_ShouldReturnPartialError() {
	transferMock := entity.NewTransfer(inputTransferMock.Payer, inputTransferMock.Payee, inputTransferMock.Amount)

	transferAuthorizationServiceMock := &mocks.TransferAuthorizationServiceMock{}
	transferAuthorizationServiceMock.On("Check", context.Background(), *transferMock).Return(nil)

	payerMock := entity.NewUser("any_first_name", "any_last_name", "any_email", "any_password", "any_taxpayer_identification", "any_role")
	payerMock.UserID = "any_payer_id"
	payeeMock := entity.NewUser("any_first_name", "any_last_name", "any_email", "any_password", "any_taxpayer_identification", "any_role")
	payeeMock.UserID = "any_payee_id"

	emailNotificationServiceMock := &mocks.EmailNotificationServiceMock{}
	emailNotificationServiceMock.On("TransferNotification", context.Background(), *payerMock, *payeeMock, *transferMock).Return(errors.New("any_error"))

	walletRepositoryMock := &mocks.WalletRepositoryMock{}
	walletRepositoryMock.On("Load", context.Background(), "any_payer_id").Return(entity.NewWallet("any_payer_id", 10), nil)
	walletRepositoryMock.On("Transfer", context.Background(), *transferMock).Return(payerMock, payeeMock, nil)

	sut := NewTransfer(transferAuthorizationServiceMock, emailNotificationServiceMock, walletRepositoryMock)
	output, err := sut.Execute(context.Background(), inputTransferMock)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), OuputTransferDTO{
		StatusCode: 206,
		Success:    true,
		Data:       "transfer successful but something went wrong to notify users",
	}, *output)

	transferAuthorizationServiceMock.AssertExpectations(s.T())
	emailNotificationServiceMock.AssertExpectations(s.T())
	walletRepositoryMock.AssertExpectations(s.T())
}

func TestTransferSuiteSuite(t *testing.T) {
	suite.Run(t, new(TransferSuite))
}
