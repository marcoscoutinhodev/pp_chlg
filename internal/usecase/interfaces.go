package usecase

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type IdentityManagerInterface interface {
	GetGroupID(ctx context.Context, group string) (groupID string, err error)
	CreateUser(ctx context.Context, user entity.User, groupID string) (*gocloak.User, error)
	AuthenticateClient(ctx context.Context, username, password string) (*gocloak.JWT, error)
}

type UserValidatorInterface interface {
	Validate(user entity.User) (errors []string)
}

type UserRepositoryInterface interface {
	CheckUserIsRegistered(ctx context.Context, user entity.User) (bool, error)
	Save(ctx context.Context, user entity.User) error
}
