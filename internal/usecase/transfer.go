package usecase

import (
	"context"
	"net/http"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type Transfer struct {
	transferAuthorizationService TransferAuthorizationServiceInterface
	emailNotificationService     EmailNotificationServiceInterface
	walletRepository             WalletRepositoryInterface
	transferRepository           TransferRepositoryInterface
}

func NewTransfer(ats TransferAuthorizationServiceInterface, ens EmailNotificationServiceInterface, wr WalletRepositoryInterface, tr TransferRepositoryInterface) *Transfer {
	return &Transfer{
		transferAuthorizationService: ats,
		emailNotificationService:     ens,
		walletRepository:             wr,
		transferRepository:           tr,
	}
}

type TransferOuputDTO struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

type TransferInputDTO struct {
	Payer string  `json:"-"`
	Payee string  `json:"payee"`
	Value float64 `json:"value"`
}

func (t Transfer) Transfer(ctx context.Context, input *TransferInputDTO) (*TransferOuputDTO, error) {
	if input.Value <= 0 {
		output := &TransferOuputDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     []string{"invalid value for transfer"},
		}
		return output, nil
	}

	transfer := entity.NewTransfer(input.Payer, input.Payee, input.Value)

	wallet, err := t.walletRepository.Load(ctx, transfer.Payer)
	if err != nil {
		return nil, err
	}

	if (wallet.Balance - input.Value) < 0 {
		output := &TransferOuputDTO{
			StatusCode: http.StatusPaymentRequired,
			Success:    false,
			Errors:     []string{"insufficient funds"},
		}
		return output, nil
	}

	err = t.transferAuthorizationService.Check(ctx, *transfer)
	if err != nil {
		output := &TransferOuputDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Success:    false,
			Errors:     []string{"the transfer was not authorized"},
		}
		return output, nil
	}

	userPayer, userPayee, err := t.walletRepository.Transfer(ctx, transfer)
	if err != nil {
		output := &TransferOuputDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Success:    false,
			Errors:     []string{"unexpected error completing transfer"},
		}
		return output, nil
	}

	if err := t.emailNotificationService.TransferNotification(ctx, *userPayer, *userPayee, *transfer); err != nil {
		return &TransferOuputDTO{
			StatusCode: http.StatusPartialContent,
			Success:    true,
			Data:       "transfer successful but something went wrong to notify users",
		}, nil
	}

	return &TransferOuputDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       "successful transfer",
	}, nil
}

type TransferListInputDTO struct {
	UserID string
	Page   int64
	Limit  int64
}

func (t Transfer) List(ctx context.Context, input *TransferListInputDTO) (*TransferOuputDTO, error) {
	if input.Page <= 0 || input.Limit <= 0 {
		return &TransferOuputDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     []string{"limit and page must be greater than 0"},
		}, nil
	}

	var transfers []entity.Transfer
	if err := t.transferRepository.List(ctx, input.UserID, input.Page, input.Limit, &transfers); err != nil {
		return nil, err
	}

	return &TransferOuputDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       transfers,
	}, nil
}
