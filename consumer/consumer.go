package consumer

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Consume() {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	// Declare DLX and DLQ
	err = ch.ExchangeDeclare("my-dlx", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare("my-dl-queue", true, false, false, false, amqp.Table{
		"x-message-ttl":             int32(5000), // 5 seconds
		"x-dead-letter-exchange":    "service_a_inner_exch",
		"x-dead-letter-routing-key": "",
	})
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind("my-dl-queue", "dlx-key", "my-dlx", false, nil)
	if err != nil {
		panic(err)
	}

	// Declare main exchange and queue
	err = ch.ExchangeDeclare("service_a_inner_exch", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare("service_a_input_q", true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    "my-dlx",
		"x-dead-letter-routing-key": "dlx-key",
	})
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind("service_a_input_q", "", "service_a_inner_exch", false, nil)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume("service_a_input_q", "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for m := range msgs {
			log.Printf("[%s] Main queue got: %s → rejecting\n", time.Now().Format("15:04:05"), m.Body)
			m.Reject(false)
		}

		log.Println()
	}()
}
