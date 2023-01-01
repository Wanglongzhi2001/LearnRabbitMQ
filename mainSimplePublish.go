package main

import (
	"fmt"
	"learn-rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("testSimple")
	defer rabbitmq.Destroy()
	rabbitmq.PublishSimple("Hello World!")
	fmt.Println("发送成功！")

}
