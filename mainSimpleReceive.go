package main

import "learn-rabbitmq/RabbitMQ"

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("testSimple") // 名字必须一样
	defer rabbitmq.Destroy()
	rabbitmq.ConsumeSimple()
}
