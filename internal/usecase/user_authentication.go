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

type OutputUserAuthenticationDTO struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

type UserAuthentication_CreateUserInputDTO struct {
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	TaxpayeerIdentification string `json:"taxpayeer_identification"`
	Role                    string `json:"role"`
}

func (u UserAuthentication) CreateUser(ctx context.Context, input *UserAuthentication_CreateUserInputDTO) (*OutputUserAuthenticationDTO, error) {
	user := entity.NewUser(input.FirstName, input.LastName, input.Email, input.Password, input.TaxpayeerIdentification, input.Role)

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
			Errors:     []string{"email and/or taxpayeer identification are already registered"},
		}
		return output, nil
	}

	kcUser, err := u.IdentityManager.CreateUser(ctx, *user)
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

type UserAuthentication_AuthenticateUserInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserAuthentication) AuthenticateUser(ctx context.Context, input *UserAuthentication_AuthenticateUserInputDTO) (*OutputUserAuthenticationDTO, error) {
	err := u.UserValidator.ValidateEmailAndPassword(input.Email, input.Password)
	if len(err) == 0 {
		jwt, err := u.IdentityManager.AuthenticateUser(ctx, input.Email, input.Password)
		if err == nil {
			output := &OutputUserAuthenticationDTO{
				StatusCode: http.StatusOK,
				Success:    true,
				Data:       jwt,
			}
			return output, nil
		}
	}

	output := &OutputUserAuthenticationDTO{
		StatusCode: http.StatusUnauthorized,
		Success:    false,
		Errors:     []string{"email and/or password are invalid"},
	}
	return output, nil
}
