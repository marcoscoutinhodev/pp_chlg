package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/http/controller"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/keycloak"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/repository"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/validator"
	"github.com/marcoscoutinhodev/pp_chlg/internal/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	identityManager := keycloak.NewIdentityManager()

	mongoClient := repository.NewMongoClient(context.Background())
	userRepository := repository.NewUserRepository(mongoClient)

	userValidator := validator.NewUserValidator()

	userAuthenticationUseCase := usecase.NewUserAuthentication(identityManager, userValidator, userRepository)
	userAuthenticationController := controller.NewUserAuthenticationController(*userAuthenticationUseCase)

	mux.Route("/user", func(r chi.Router) {
		r.Post("/signup", userAuthenticationController.CreateUser)
		r.Post("/signin", userAuthenticationController.AuthenticateUser)
	})

	http.ListenAndServe(os.Getenv("SERVER_PORT"), mux)
}
