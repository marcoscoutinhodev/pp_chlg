package usecase

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type IdentityManagerInterface interface {
	CreateUser(ctx context.Context, user entity.User) (*gocloak.User, error)
	AuthenticateUser(ctx context.Context, username, password string) (*gocloak.JWT, error)
}

type UserValidatorInterface interface {
	Validate(user entity.User) (errors []string)
	ValidateEmailAndPassword(email, password string) (errors []string)
}

type UserRepositoryInterface interface {
	CheckUserIsRegistered(ctx context.Context, user entity.User) (bool, error)
	Save(ctx context.Context, user entity.User) error
}

type WalletRepositoryInterface interface {
	Load(ctx context.Context, userID string) (*entity.Wallet, error)
	Transfer(ctx context.Context, transfer *entity.Transfer) (userPayer, userPayee *entity.User, err error)
}

type TransferRepositoryInterface interface {
	List(ctx context.Context, userID string, page, limit int64, transfers *[]entity.Transfer) error
}
type TransferAuthorizationServiceInterface interface {
	Check(ctx context.Context, transfer entity.Transfer) error
}

type EmailNotificationServiceInterface interface {
	TransferNotification(ctx context.Context, payer, payee entity.User, transfer entity.Transfer) error
}
