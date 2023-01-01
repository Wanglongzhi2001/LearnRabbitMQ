package main

import (
	"fmt"
	"learn-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" + "testSimple")
	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello test!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
