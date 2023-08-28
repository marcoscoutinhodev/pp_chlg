package mocks

import (
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type EmailNotificationServiceMock struct {
	mock.Mock
}

func (e *EmailNotificationServiceMock) TransferNotification(payer, payee *entity.User) {
	e.Called(payer, payee)
}
