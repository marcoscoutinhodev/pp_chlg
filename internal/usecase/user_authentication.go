package usecase

import (
	"context"
	"net/http"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type UserAuthentication struct {
	IdentityManager IdentityManagerInterface
	UserValidator   UserValidatorInterface
	UserRepository  UserRepositoryInterface
}

func NewUserAuthentication(im IdentityManagerInterface, uv UserValidatorInterface, ur UserRepositoryInterface) *UserAuthentication {
	return &UserAuthentication{
		IdentityManager: im,
		UserValidator:   uv,
		UserRepository:  ur,
	}
}

type UserAuthentication_CreateUserInputDTO struct {
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	TaxpayeerIdentification string `json:"taxpayeer_identification"`
	Group                   string `json:"group"`
}

type OutputUserAuthenticationDTO struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

func (u UserAuthentication) CreateUser(ctx context.Context, input *UserAuthentication_CreateUserInputDTO) (*OutputUserAuthenticationDTO, error) {
	user := entity.NewUser(input.FirstName, input.LastName, input.Email, input.Password, input.TaxpayeerIdentification, input.Group)

	if err := u.UserValidator.Validate(*user); len(err) > 0 {
		output := &OutputUserAuthenticationDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     err,
		}
		return output, nil
	}

	userIsRegistered, err := u.UserRepository.CheckUserIsRegistered(ctx, *user)
	if err != nil {
		return nil, err
	}

	if userIsRegistered {
		output := &OutputUserAuthenticationDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     []string{"email or and taxpayeer identification is already registered"},
		}
		return output, nil
	}

	groupID, err := u.IdentityManager.GetGroupID(ctx, user.Group)
	if err != nil {
		return nil, err
	}

	if groupID == "" {
		output := &OutputUserAuthenticationDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     []string{"invalid group provided"},
		}
		return output, nil
	}

	kcUser, err := u.IdentityManager.CreateUser(ctx, *user, groupID)
	if err != nil {
		return nil, err
	}

	user.UserID = *kcUser.ID

	if u.UserRepository.Save(ctx, *user); err != nil {
		return nil, err
	}

	return &OutputUserAuthenticationDTO{
		StatusCode: http.StatusCreated,
		Success:    true,
	}, nil
}
