package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	emailsendingsimulation "github.com/marcoscoutinhodev/pp_chlg/email-sending-simulation"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/http/controller"
	http_middleware "github.com/marcoscoutinhodev/pp_chlg/internal/infra/http/middleware"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/keycloak"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/repository"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/service"
	"github.com/marcoscoutinhodev/pp_chlg/internal/infra/validator"
	"github.com/marcoscoutinhodev/pp_chlg/internal/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// essa goroutine vai emitir logs
	// simulando algum serviço que consuma os dados de transferência
	// para fazer o envio dos emails
	go emailsendingsimulation.EmailSendingSimulation()
}

func main() {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	identityManager := keycloak.NewIdentityManager()

	// repositories
	mongoClient := repository.NewMongoClient(context.Background())
	userRepository := repository.NewUserRepository(mongoClient)
	walletRepository := repository.NewWalleteRepository(mongoClient)
	transferRepository := repository.NewTransferRepository(mongoClient)

	// services
	transferAuthorizationService := service.NewTransferAuthorizationService()
	emailNotificationService := service.NewEmailNotificationService()

	// validators
	userValidator := validator.NewUserValidator()

	// usecases
	userAuthenticationUseCase := usecase.NewUserAuthentication(identityManager, userValidator, userRepository)
	transferUseCase := usecase.NewTransfer(transferAuthorizationService, emailNotificationService, walletRepository, transferRepository)

	// controllers
	userAuthenticationController := controller.NewUserAuthenticationController(*userAuthenticationUseCase)
	transferController := controller.NewTransfer(*transferUseCase)

	// middlewares
	transferListAuthorizationMiddleware := http_middleware.NewAuthorization(map[string]bool{"shopkeeper": true, "customer": true})
	transferAuthorizationMiddleware := http_middleware.NewAuthorization(map[string]bool{"customer": true})

	mux.Route("/user", func(r chi.Router) {
		r.Post("/signup", userAuthenticationController.CreateUser)
		r.Post("/signin", userAuthenticationController.AuthenticateUser)
	})

	mux.Route("/transfer", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			transferAuthorizationMiddleware.Handle(w, r, transferController.Transfer)
		})
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			transferListAuthorizationMiddleware.Handle(w, r, transferController.List)
		})
	})

	http.ListenAndServe(os.Getenv("SERVER_PORT"), mux)
}
