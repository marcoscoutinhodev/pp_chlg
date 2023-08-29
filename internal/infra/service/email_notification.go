package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/marcoscoutinhodev/pp_chlg/internal/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailNotificationService struct {
	conn *amqp.Connection
}

func NewEmailNotificationService() *EmailNotificationService {
	c, err := amqp.Dial(os.Getenv("RBMQ_URI"))
	if err != nil {
		panic(err)
	}

	return &EmailNotificationService{
		conn: c,
	}
}

type transferNotificationdDTO struct {
	PayerEmail        string  `json:"payer_email"`
	PayerName         string  `json:"payer_name"`
	PayeeEmail        string  `json:"payee_email"`
	PayeeName         string  `json:"payee_name"`
	AmountTransferred float64 `json:"amount_transferred"`
	Date              string  `json:"date"`
}

func (e EmailNotificationService) TransferNotification(ctx context.Context, payer, payee entity.User, transfer entity.Transfer) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	input := transferNotificationdDTO{
		PayerEmail:        payer.Email,
		PayerName:         fmt.Sprintf("%s %s", payer.FirstName, payer.LastName),
		PayeeEmail:        payee.Email,
		PayeeName:         fmt.Sprintf("%s %s", payee.FirstName, payee.LastName),
		AmountTransferred: transfer.Amount,
		Date:              transfer.Date,
	}
	body, err := json.Marshal(&input)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		os.Getenv("RBMQ_TRANSFER_NOTIFICATION_EXCHANGE"),
		os.Getenv("RBMQ_TRANSFER_NOTIFICATION_KEY"),
		true,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
