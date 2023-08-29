package mocks

import (
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type EmailNotificationServiceMock struct {
	mock.Mock
}

func (e *EmailNotificationServiceMock) TransferNotification(payer, payee *entity.User, amount float64) {
	e.Called(payer, payee, amount)
}
