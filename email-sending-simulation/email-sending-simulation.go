package emailsendingsimulation

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// essa função serve apenas para emitir logs
// simulando algum serviço que consuma os dados de transferência
// para fazer o envio dos emails
func EmailSendingSimulation() {
	var attempts uint
	var err error
	if err := emailSendingSimulation(attempts, err); err != nil {
		fmt.Printf("EmailSendingSimulation: %v\n", err)
	}
}

func emailSendingSimulation(attempts uint, err error) error {
	if attempts < 10 {
		c, err := amqp.Dial(os.Getenv("RBMQ_URI"))
		if err != nil {
			time.Sleep(time.Minute)
			return emailSendingSimulation(attempts+1, err)
		}
		defer c.Close()

		ch, err := c.Channel()
		if err != nil {
			time.Sleep(time.Minute)
			return emailSendingSimulation(attempts+1, err)
		}
		defer ch.Close()

		delivery, err := ch.Consume(
			os.Getenv("RBMQ_TRANSFER_NOTIFICATION_QUEUE"),
			*new(string),
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			time.Sleep(time.Minute)
			return emailSendingSimulation(attempts+1, err)
		}

		transferNotification := struct {
			TransferID       string  `json:"transfer_id"`
			PayerEmail       string  `json:"payer_email"`
			PayerName        string  `json:"payer_name"`
			PayeeEmail       string  `json:"payee_email"`
			PayeeName        string  `json:"payee_name"`
			ValueTransferred float64 `json:"value_transferred"`
			Date             string  `json:"date"`
		}{}

		for d := range delivery {
			if err := json.Unmarshal(d.Body, &transferNotification); err == nil {
				fmt.Println("*****************************************************************************************")

				fmt.Printf("Enviando email para %s\nTitulo: Pagamento Realizado\nMensagem: Você pagou R$%.2f a %s | Transação: %s | %s\n\n\n",
					transferNotification.PayerEmail, transferNotification.ValueTransferred, transferNotification.PayeeName, transferNotification.TransferID, transferNotification.Date)

				fmt.Printf("Enviando email para %s\nTitulo: Pagamento Recebido\nMensagem: %s, você recebeu R$%.2f | Transação: %s | %s\n",
					transferNotification.PayeeEmail, transferNotification.PayeeName, transferNotification.ValueTransferred, transferNotification.TransferID, transferNotification.Date)

				fmt.Printf("*****************************************************************************************\n\n")

				ch.Ack(d.DeliveryTag, true)
			} else {
				ch.Nack(d.DeliveryTag, false, true)
			}
		}
	}

	return err
}
