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

type UserAuthenticationOutputDTO struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
}

type UserCreateInputDTO struct {
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Email                  string `json:"email"`
	Password               string `json:"password"`
	TaxpayerIdentification string `json:"taxpayer_identification"`
	Role                   string `json:"role"`
}

func (u UserAuthentication) CreateUser(ctx context.Context, input *UserCreateInputDTO) (*UserAuthenticationOutputDTO, error) {
	user := entity.NewUser(input.FirstName, input.LastName, input.Email, input.Password, input.TaxpayerIdentification, input.Role)

	if err := u.UserValidator.Validate(*user); len(err) > 0 {
		output := &UserAuthenticationOutputDTO{
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
		output := &UserAuthenticationOutputDTO{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Errors:     []string{"email and/or taxpayer identification are already registered"},
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

	return &UserAuthenticationOutputDTO{
		StatusCode: http.StatusCreated,
		Success:    true,
	}, nil
}

type UserAuthenticateInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserAuthentication) AuthenticateUser(ctx context.Context, input *UserAuthenticateInputDTO) (*UserAuthenticationOutputDTO, error) {
	err := u.UserValidator.ValidateEmailAndPassword(input.Email, input.Password)
	if len(err) == 0 {
		jwt, err := u.IdentityManager.AuthenticateUser(ctx, input.Email, input.Password)
		if err == nil {
			output := &UserAuthenticationOutputDTO{
				StatusCode: http.StatusOK,
				Success:    true,
				Data:       jwt,
			}
			return output, nil
		}
	}

	output := &UserAuthenticationOutputDTO{
		StatusCode: http.StatusUnauthorized,
		Success:    false,
		Errors:     []string{"email and/or password are invalid"},
	}
	return output, nil
}
