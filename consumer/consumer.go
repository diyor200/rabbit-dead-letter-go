package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume() {
	conn, _ := amqp.Dial("amqp://user:password@localhost:5672")
	ch, _ := conn.Channel()

	// DLX
	ch.ExchangeDeclare("my-dlx", "topic", true, false, false, false, nil)
	ch.QueueDeclare("my-dl-queue", true, false, false, false, nil)
	ch.QueueBind("my-dl-queue", "dlx-key", "my-dlx", false, nil)

	// Main exchange
	ch.ExchangeDeclare("service_a_inner_exch", "topic", true, false, false, false, nil)

	args := amqp.Table{
		"x-dead-letter-exchange":    "my-dlx",
		"x-dead-letter-routing-key": "dlx-key",
	}

	ch.QueueDeclare("service_a_input_q", true, false, false, false, args)
	ch.QueueBind("service_a_input_q", "", "service_a_inner_exch", false, nil)

	msgs, _ := ch.Consume("service_a_input_q", "", false, false, false, false, nil)

	go func() {
		for m := range msgs {
			log.Printf("Main queue got: %s\n", m.Body)

			m.Reject(false)
		}
	}()
}

func ConsumeDlQ() {
	conn, _ := amqp.Dial("amqp://user:password@localhost:5672")
	ch, _ := conn.Channel()

	ch.QueueDeclare("my-dl-queue", true, false, false, false, nil)
	msgs, _ := ch.Consume("my-dl-queue", "", true, false, false, false, nil)

	go func() {
		for m := range msgs {
			log.Printf("DLQ got: %s", m.Body)
		}
	}()
}
