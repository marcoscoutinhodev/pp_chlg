package mocks

import (
	"context"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) CheckUserIsRegistered(ctx context.Context, user entity.User) (bool, error) {
	args := u.Called(ctx, user)
	return args.Bool(0), args.Error(1)
}

func (u *UserRepositoryMock) Save(ctx context.Context, user entity.User) error {
	args := u.Called(ctx, user)
	return args.Error(0)
}
