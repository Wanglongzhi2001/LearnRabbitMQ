package main

import (
	"fmt"
	"learn-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	Routing1 := RabbitMQ.NewRabbitMQRouting("exRouting", "Routing1")
	Routing2 := RabbitMQ.NewRabbitMQRouting("exRouting", "Routing2")
	for i := 0; i <= 10; i++ {
		Routing1.PublishRouting("Hello Routing1! " + strconv.Itoa(i))
		Routing2.PublishRouting("Hello Routing2! " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
