package service

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBrokerService struct {
	Ch *amqp.Channel
}

func NewMessageBrokerService(ch *amqp.Channel) MessageBrokerService {
	return MessageBrokerService{
		Ch: ch,
	}
}

func (m MessageBrokerService) PublishMessage(message string) error {
	q, err := m.Ch.QueueDeclare(
		"ungraded_5", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	ctx := context.TODO()

	err = m.Ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	log.Printf(" [x] Sent %s\n", message)
	return nil
}
