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
}

func NewTransfer(ats TransferAuthorizationServiceInterface, ens EmailNotificationServiceInterface, wr WalletRepositoryInterface) *Transfer {
	return &Transfer{
		transferAuthorizationService: ats,
		emailNotificationService:     ens,
		walletRepository:             wr,
	}
}

type InputTransferDTO struct {
	Payer  string  `json:"-"`
	Payee  string  `json:"payee"`
	Amount float64 `json:"amount"`
}

type OuputTransferDTO struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

func (t Transfer) Execute(ctx context.Context, input *InputTransferDTO) (*OuputTransferDTO, error) {
	transfer := entity.NewTransfer(input.Payer, input.Payee, input.Amount)

	wallet, err := t.walletRepository.Load(ctx, transfer.Payer)
	if err != nil {
		return nil, err
	}

	if (wallet.Balance - input.Amount) < 0 {
		output := &OuputTransferDTO{
			StatusCode: http.StatusPaymentRequired,
			Success:    false,
			Errors:     []string{"insufficient funds"},
		}
		return output, nil
	}

	err = t.transferAuthorizationService.Check(ctx, *transfer)
	if err != nil {
		output := &OuputTransferDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Success:    false,
			Errors:     []string{"the transfer was not authorized"},
		}
		return output, nil
	}

	userPayer, userPayee, err := t.walletRepository.Transfer(ctx, *transfer)
	if err != nil {
		output := &OuputTransferDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Success:    false,
			Errors:     []string{"unexpected error completing transfer"},
		}
		return output, nil
	}

	t.emailNotificationService.TransferNotification(userPayer, userPayee, transfer.Amount)

	return &OuputTransferDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       "successful transfer",
	}, nil
}
