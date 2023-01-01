package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" + "testSimple")
	rabbitmq.ConsumeSimple()
}
