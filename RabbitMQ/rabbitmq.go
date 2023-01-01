package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url 格式 amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://testuser:testuser@127.0.0.1:5672/test"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key string
	// 连接信息
	Mqurl string
}

// 构造函数创建RabbitMQ实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ { // 根据这三个参数的不同组合实现不同的模式
	rabbitmq := RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	// 创建RabbitMQ连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "获取channel失败")
	return &rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.conn.Close()
	r.channel.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Simple模式step1:在Simple模式下创建RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "") // exchange使用默认的(direct),key为空(不同模式传入不同参数)
}

// Simple模式step2:在Simple模式下生产代码
func (r *RabbitMQ) PublishSimple(msg string) {
	// 1.声明队列
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除(当最后一个订阅者取消订阅时)
		false,       // 是否具有排他性
		false,       // 是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 如果为true，根据exchange类型和routekey规则，如果找不到符合条件的队列那么会把发送的消息返回给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

}

// 消费消息
func (r *RabbitMQ) ConsumeSimple() {
	// 1.声明队列
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除(当最后一个订阅者取消订阅时)
		false,       // 是否具有排他性
		false,       // 是否阻塞
		nil,         //额外属性
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列是否阻塞
		false,
		// 其他参数
		nil)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	// 启用协程处理消息(因为我们生产消息是异步的)
	go func() {
		for d := range msgs {
			// 实现我们的逻辑函数
			log.Printf("received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	// 使之阻塞
	<-forever
}

// 订阅模式创建rabbitmq实例(需要设置交换机名称)
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel")
	return rabbitmq
}

// 订阅模式下生产
func (r *RabbitMQ) PublishPub(msg string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // 广播
		true,
		false,
		false, //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnError(err, "Failed to declare an exchange")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"", // 对于fanout的交换机，key为空
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.failOnError(err, "Failed to declare an publish")
}

// 订阅模式消费端代码
func (r *RabbitMQ) ReceiveSub() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "failed to declare a exchange")

	// 2.尝试创建队列，这里注意队列名称不要写(随机创建队列)
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err, "failed to declare a queue")
	// 绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		"", //在pub/sub模式下，这里的key要为空
		r.Exchange,
		false,
		nil)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("To exit press CTRL+C\n")
	<-forever
}

// 路由模式
// 创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	// 创建rabbitmq实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "failed to connect to rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel")
	return rabbitmq
}

// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(msg string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", // 要改成direct
		true,
		false,
		false,
		false,
		nil)
	r.failOnError(err, "failed to declare a exchange")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, // 要设置
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.failOnError(err, "Failed to declare an publish")
}

// 路由模式接受消息
func (r *RabbitMQ) ReceiveRouting() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	r.failOnError(err, "failed to declare an exchange")
	// 2.尝试创建队列，这里注意队列名字不要写
	q, err := r.channel.QueueDeclare(
		"", //随机产生队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	// 3.绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	// 4.消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("To exit press CTRL+C")
	<-forever
}

// 话题模式
// 创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnError(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnError(err, "failed to open a channel!")
	return rabbitmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublishTopic(msg string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err, "failed to declare an exchange")
	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, // 要设置
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	r.failOnError(err, "Failed to declare an publish")
}

// 话题模式接收消息
// 要注意key，规则
// 其中"*"用于匹配一个单词，"#"用于匹配多个单词(可以是0个)
// 匹配imooc.* 表示匹配 imooc.hello，但是imooc.hello.one需要用imooc.#才能匹配到
func (r *RabbitMQ) ReceiveTopic() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	r.failOnError(err, "failed to declare an exchange")
	// 2.尝试创建队列，这里注意队列名字不要写
	q, err := r.channel.QueueDeclare(
		"", //随机产生队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	// 3.绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	// 4.消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("To exit press CTRL+C")
	<-forever
}
