package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" + "newProduct")
	rabbitmq.ReceiveSub()
}
