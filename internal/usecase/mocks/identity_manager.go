package mocks

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/stretchr/testify/mock"
)

type IdentityManagerMock struct {
	mock.Mock
}

func (i *IdentityManagerMock) CreateUser(ctx context.Context, user entity.User) (*gocloak.User, error) {
	args := i.Called(ctx, user)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gocloak.User), args.Error(1)
}

func (i *IdentityManagerMock) AuthenticateUser(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	args := i.Called(ctx, username, password)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gocloak.JWT), args.Error(1)
}
