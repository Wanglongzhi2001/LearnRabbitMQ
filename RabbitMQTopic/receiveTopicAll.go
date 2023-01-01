package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	Topic1 := RabbitMQ.NewRabbitMQTopic("exTopic", "#")
	Topic1.ReceiveTopic()
}
