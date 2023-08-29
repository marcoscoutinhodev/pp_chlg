package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/marcoscoutinhodev/pp_chlg/internal/usecase"
)

type UserAuthenticationController struct {
	UserAuthenticationUseCase usecase.UserAuthentication
}

func NewUserAuthenticationController(uauc usecase.UserAuthentication) *UserAuthenticationController {
	return &UserAuthenticationController{
		UserAuthenticationUseCase: uauc,
	}
}

func (u UserAuthenticationController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.UserCreateInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "failed to parse request body",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	output, err := u.UserAuthenticationUseCase.CreateUser(ctx, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Printf("internal error: %v\n", err)
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}

func (u UserAuthenticationController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var input usecase.UserAuthenticateInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "failed to parse request body",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*3)
	defer cancel()

	output, err := u.UserAuthenticationUseCase.AuthenticateUser(ctx, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal error, try again in a few minutes",
		})
		fmt.Printf("internal error: %v\n", err)
		return
	}

	w.WriteHeader(output.StatusCode)
	json.NewEncoder(w).Encode(output)
}
