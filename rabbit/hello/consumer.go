package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// 1.建立连接
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Printf("connect rabbitmq failed: %s\n", err)
	}
	defer conn.Close()
	// 2.获取通道
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("create channel failed: %s\n", err)
	}
	defer ch.Close()
	// 3.声明队列
	queue, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Printf("declare queue failed: %s\n", err)
	}
	// 4.获取消息
	deliver, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("register consumer failed: %s\n", err)
	}
	go func() {
		for msg := range deliver {
			log.Printf("receive msg: %s\n", msg.Body)
		}
	}()
	log.Println("waiting for msg....")
	<-deliver
}
