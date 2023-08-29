package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type EmailNotificationServiceMock struct {
	mock.Mock
}

func (e *EmailNotificationServiceMock) TransferNotification(ctx context.Context, payer, payee entity.User, transfer entity.Transfer) error {
	args := e.Called(ctx, payer, payee, transfer)
	return args.Error(0)
}
