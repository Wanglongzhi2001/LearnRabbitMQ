package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	Topic1 := RabbitMQ.NewRabbitMQTopic("exTopic", "test.*.two")
	Topic1.ReceiveTopic()
}
