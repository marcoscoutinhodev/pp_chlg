package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type TransferAuthorizationService struct{}

func NewTransferAuthorizationService() *TransferAuthorizationService {
	return &TransferAuthorizationService{}
}

func (t TransferAuthorizationService) Check(ctx context.Context, transfer entity.Transfer) error {
	res, err := http.Get(os.Getenv("TRANSFER_AUTHORIZATION_URL"))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("unauthorized transfer")
	}

	result := struct {
		Message string `json:"message"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	if result.Message != "Autorizado" {
		return errors.New("unauthorized transfer")
	}

	return nil
}
