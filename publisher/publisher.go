package publisher

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(message string) {
	conn, _ := amqp.Dial("amqp://user:password@localhost:5672")
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	ch.ExchangeDeclare("service_a_inner_exch", "topic", true, false, false, false, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		"service_a_inner_exch",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Headers:     amqp.Table{"x-retries": int32(0)},
		},
	)
	if err != nil {
		log.Println("Publish failed:", err)
	}
}
