package main

import (
	"rabbit_dead_letter/consumer"
	"rabbit_dead_letter/publisher"
	"time"
)

func main() {
	go consumer.Consume()
	go consumer.ConsumeDlQ()

	time.Sleep(time.Second * 2) // wait for consumers

	for {
		publisher.Publish()
		time.Sleep(time.Second)
	}
}
