package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type TransferAuthorizationServiceMock struct {
	mock.Mock
}

func (t *TransferAuthorizationServiceMock) Check(ctx context.Context, transfer entity.Transfer) error {
	args := t.Called(ctx, transfer)
	return args.Error(0)
}
