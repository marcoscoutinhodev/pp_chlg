package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type WalletRepositoryMock struct {
	mock.Mock
}

func (w *WalletRepositoryMock) Load(ctx context.Context, userID string) (*entity.Wallet, error) {
	args := w.Called(ctx, userID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (w *WalletRepositoryMock) Transfer(ctx context.Context, transfer entity.Transfer) (userPayer, userPayee *entity.User, err error) {
	args := w.Called(ctx, transfer)

	if userPayer == nil || userPayee == nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(*entity.User), args.Get(1).(*entity.User), args.Error(2)
}
