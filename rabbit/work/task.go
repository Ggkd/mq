package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

//func main() {
//	// 1.建立连接
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
//	if err != nil {
//		log.Printf("connect rabbitmq failed: %s\n", err)
//	}
//	defer conn.Close()
//	// 2.创建通道
//	ch, err := conn.Channel()
//	if err != nil {
//		log.Printf("create channel failed: %s\n", err)
//	}
//	defer ch.Close()
//	// 3.声明要送消息的队列
//	queue, err := ch.QueueDeclare("hello",
//		false,
//		false,
//		false,
//		false,
//		nil)
//	if err != nil {
//		log.Printf("declare queue failed: %s\n", err)
//	}
//	// 4.发送消息到队列
//	body := bodyFrom(os.Args)
//	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
//		DeliveryMode: amqp.Persistent,  // 持久化
//		Body:         []byte(body),
//	})
//	if err != nil {
//		log.Printf("publish msg failed\": %s\n", err)
//	}
//}

func main() {
	// 1.建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Printf("connect rabbitmq failed: %s\n", err)
	}
	defer conn.Close()
	// 2.创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("create channel failed: %s\n", err)
	}
	defer ch.Close()
	// 3.声明要送消息的队列
	queue, err := ch.QueueDeclare("task_work",
		true, // 持久化队列
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Printf("declare queue failed: %s\n", err)
	}
	// 4.发送消息到队列
	body := bodyFrom(os.Args)
	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent, // 持久化
		Body:         []byte(body),
	})
	if err != nil {
		log.Printf("publish msg failed\": %s\n", err)
	}
}
