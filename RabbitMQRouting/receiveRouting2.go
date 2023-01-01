package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	Routing2 := RabbitMQ.NewRabbitMQRouting("exRouting", "Routing2")
	Routing2.ReceiveRouting()
}
