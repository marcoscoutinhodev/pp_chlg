package http_middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserIDKeyContext struct{}

type Authorization struct {
	Roles map[string]bool
}

func NewAuthorization(roles map[string]bool) *Authorization {
	return &Authorization{
		Roles: roles,
	}
}

func (a Authorization) Handle(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	accessToken := r.Header.Get("x_access_token")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	parts := strings.Split(accessToken, " ")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	if parts[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("KC_TOKEN_PUBLIC_KEY")))
	if err != nil {
		fmt.Printf("internal error: fail to parse public key (%v)\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "internal server error",
		})
		return
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	resourceAccess := claims["resource_access"].(map[string]interface{})
	resourceAccessAll := resourceAccess["all"].(map[string]interface{})
	roles := resourceAccessAll["roles"].([]interface{})
	for _, role := range roles {
		if !a.Roles[role.(string)] {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "forbidden",
			})
			return
		}
	}

	request := r.WithContext(context.WithValue(r.Context(), UserIDKeyContext{}, claims["sub"]))
	next(w, request)
}
