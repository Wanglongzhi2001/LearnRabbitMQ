package main

import (
	"fmt"
	"learn-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	Topic1 := RabbitMQ.NewRabbitMQTopic("exTopic", "test.topic.one")
	Topic2 := RabbitMQ.NewRabbitMQTopic("exTopic", "test.topic.two")
	for i := 0; i <= 10; i++ {
		Topic1.PublishTopic("Hello Topic1! " + strconv.Itoa(i))
		Topic2.PublishTopic("Hello Topic2! " + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
