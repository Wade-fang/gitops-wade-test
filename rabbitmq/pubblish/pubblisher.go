package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:12345fw@136.116.224.158:5672/")
	failOnError(err, "failed to connect rabbitmq server")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to create channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"log",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to create fanout exchange")

	body := "hello world"
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	err = ch.PublishWithContext(ctx, "log", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

}
