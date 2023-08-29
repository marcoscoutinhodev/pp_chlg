package service

import (
	"fmt"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
)

type EmailNotificationService struct{}

func NewEmailNotificationService() *EmailNotificationService {
	return &EmailNotificationService{}
}

func (e EmailNotificationService) TransferNotification(payer, payee *entity.User, amount float64) {
	fmt.Printf("Enviando email para cliente: %s | %s. Total transferido: %.2f\n", payer.ID, payer.Email, amount)
	fmt.Printf("Enviando email para cliente: %s | %s. Total Recebido: %.2f\n", payee.ID, payee.Email, amount)
}
