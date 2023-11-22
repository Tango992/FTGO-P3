package main

import (
	"context"
	"log"
	"rabbit/helpers"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"demo-ftgo-hello", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	ctx := context.TODO()

	body := "Foo Bar"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	helpers.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
