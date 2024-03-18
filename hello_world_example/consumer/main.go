package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Hello",
		false,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed To declared a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish message")
	log.Printf("[x] Sent %s\n", body)
}
