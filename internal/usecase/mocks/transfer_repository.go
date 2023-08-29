package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type TransferRepositoryMock struct {
	mock.Mock
}

func (w *TransferRepositoryMock) List(ctx context.Context, userID string, page, limit int64, transfers *[]entity.Transfer) error {
	args := w.Called(ctx, userID, page, limit, transfers)
	return args.Error(0)
}
