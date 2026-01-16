package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"strings"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {
	conn, err := amqp.Dial("amqp://admin:12345fw@136.116.224.158:5672/")
	failOnError(err, "无法连接到rabbitmq服务器")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "无法打开通道")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "无法声明队列")

	duration := time.Second * 5
	ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
	defer cancelFunc()

	body := bodyFrom(os.Args)
	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "推送消息失败")
	log.Printf(" [x] Sent %s\n", body)
}
