package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	Routing1 := RabbitMQ.NewRabbitMQRouting("exRouting", "Routing1")
	Routing1.ReceiveRouting()
}
