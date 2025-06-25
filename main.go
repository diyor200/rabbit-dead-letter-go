package main

import (
	"fmt"
	"rabbit_dead_letter/consumer"
	"rabbit_dead_letter/publisher"
	"time"
)

func main() {
	go consumer.Consume()
	// go consumer.ConsumeDlQ()

	time.Sleep(time.Second * 2) // wait for consumers

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("message: %d", i)
		publisher.Publish(msg)
	}
	select {}
}
