package usecase

import (
	"context"

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

type InputUserAuthenticationCreateUserDTO struct {
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	TaxpayeerIdentification string `json:"taxpayeer_identification"`
	Group                   string `json:"group"`
}

type OutputUserAuthenticationCreateUserDTO struct {
	ValidationError  []string `json:"validation_error,omitempty"`
	DuplicationError string   `json:"duplication_error,omitempty"`
}

func (u UserAuthentication) CreateUser(ctx context.Context, input *InputUserAuthenticationCreateUserDTO) (*OutputUserAuthenticationCreateUserDTO, error) {
	user := entity.NewUser(input.FirstName, input.LastName, input.Email, input.Password, input.TaxpayeerIdentification, input.Group)

	if err := u.UserValidator.Validate(*user); len(err) > 0 {
		output := &OutputUserAuthenticationCreateUserDTO{ValidationError: err}
		return output, nil
	}

	userIsRegistered, err := u.UserRepository.CheckUserIsRegistered(ctx, *user)
	if err != nil {
		return nil, err
	}

	if userIsRegistered {
		output := &OutputUserAuthenticationCreateUserDTO{DuplicationError: "email or and taxpayeer identification is already registered"}
		return output, nil
	}

	groupID, err := u.IdentityManager.GetGroupID(ctx, user.Group)
	if err != nil {
		return nil, err
	}

	if groupID == "" {
		output := &OutputUserAuthenticationCreateUserDTO{ValidationError: []string{"invalid group provided"}}
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

	return &OutputUserAuthenticationCreateUserDTO{}, nil
}
