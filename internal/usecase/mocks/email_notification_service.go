package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type EmailNotificationServiceMock struct {
	mock.Mock
}

func (e *EmailNotificationServiceMock) TransferNotification(ctx context.Context, payer, payee entity.User, amount float64) error {
	args := e.Called(ctx, payer, payee, amount)
	return args.Error(0)
}
