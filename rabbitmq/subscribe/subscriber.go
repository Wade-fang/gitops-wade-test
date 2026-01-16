package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
	failOnError(err, "failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "failed to declare queue")

	err = ch.ExchangeDeclare(
		"log",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to declare exchange")

	err = ch.QueueBind(
		q.Name,
		"",
		"log",
		false,
		nil,
	)
	failOnError(err, "failed to bind queue on exchange")

	var forever chan struct{}
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	go func() {
		for v := range msgs {
			fmt.Printf("[x] msg: %s\n", v.Body)
			v.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}
