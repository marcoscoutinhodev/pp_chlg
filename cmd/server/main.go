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

	mongoClient := repository.NewMongoClient(context.Background())

	identityManager := keycloak.NewIdentityManager()
	userValidator := validator.NewUserValidator()
	userRepository := repository.NewUserRepository(mongoClient)
	userAuthenticationUseCase := usecase.NewUserAuthentication(identityManager, userValidator, userRepository)
	userAuthenticationController := controller.NewUserAuthenticationController(*userAuthenticationUseCase)

	mux.Route("/user", func(r chi.Router) {
		r.Post("/signup", userAuthenticationController.CreateUser)
		r.Post("/signin", userAuthenticationController.AuthenticateUser)
	})

	transferAuthorizationService := service.NewTransferAuthorizationService()
	emailNotificationService := service.NewEmailNotificationService()
	walletRepository := repository.NewWalleteRepository(mongoClient)
	transferUseCase := usecase.NewTransfer(transferAuthorizationService, emailNotificationService, walletRepository)
	transferController := controller.NewTransfer(*transferUseCase)
	transferAuthorizationMiddleware := http_middleware.NewAuthorization(map[string]bool{"customer": true})

	mux.Post("/transfer", func(w http.ResponseWriter, r *http.Request) {
		transferAuthorizationMiddleware.Handle(w, r, transferController.Handle)
	})

	http.ListenAndServe(os.Getenv("SERVER_PORT"), mux)
}
