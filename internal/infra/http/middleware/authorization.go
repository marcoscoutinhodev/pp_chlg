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
	Roles []string
}

func NewAuthorization(roles []string) *Authorization {
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

	roles := claims["roles"]
	if rolesFromAccessTokenAsMap, ok := roles.(map[string]interface{}); ok {
		for _, role := range a.Roles {
			// all roles must be registered in the user's access token to continue
			if rolesFromAccessTokenAsMap[role] == nil {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": false,
					"error":   "forbidden",
				})
				return
			}
		}
	}

	request := r.WithContext(context.WithValue(r.Context(), UserIDKeyContext{}, claims["sub"]))
	next(w, request)
}
