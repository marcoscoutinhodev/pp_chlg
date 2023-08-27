package validator

import (
	"net/mail"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	"github.com/marcoscoutinhodev/pp_chlg/pkg"
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (uv UserValidator) Validate(user entity.User) (errors []string) {
	if isValid := uv.NameValidator(user.FirstName); !isValid {
		errors = append(errors, "first name must have only upper case letters")
	}

	if isValid := uv.NameValidator(user.LastName); !isValid {
		errors = append(errors, "last name must have only upper case letters")
	}

	if isValid := uv.EmailValidator(user.Email); !isValid {
		errors = append(errors, "email must be a valid email address")
	}

	if isValid := uv.PasswordValidator(user.Password); !isValid {
		errors = append(errors, "password must be at least 7 characters with uppercase/lowercase letters, number, special character")
	}

	if isValid := pkg.CPFCNPJValidator(user.TaxpayeerIdentification); !isValid {
		errors = append(errors, "taxpayeer identification must be a valid CPF or CNPJ")
	}

	if isValid := uv.GroupValidator(user.Group); !isValid {
		errors = append(errors, "unknown group")
	}

	return
}

func (uv UserValidator) NameValidator(name string) bool {
	isMatch := regexp.MustCompile(`^[A-Z\s]*$`).MatchString(name)
	return isMatch
}

func (uv UserValidator) EmailValidator(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (uv UserValidator) PasswordValidator(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 7 {
		hasMinLen = true

		for _, char := range password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func (uv UserValidator) GroupValidator(group string) bool {
	availableGroups := strings.Split(os.Getenv("KC_AVAILABLE_GROUPS"), ",")

	for _, g := range availableGroups {
		if group == g {
			return true
		}
	}

	return false
}
