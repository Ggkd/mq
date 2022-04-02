package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
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
	queue, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Printf("declare queue failed: %s\n", err)
	}
	// 3.1 预处理
	ch.Qos(2,
		0,
		false,
	)
	// 4.获取消息
	msg, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Printf("register consumer failed: %s\n", err)
	}
	go func() {
		for m := range msg {
			log.Printf("receive msg: %s\n", m.Body)
			time.Sleep(time.Second * 1)
			log.Println("Done!")
			// 回复应答
			err = ch.Publish(
				"",
				m.ReplyTo,
				false,
				false,
				amqp.Publishing{CorrelationId: m.CorrelationId, Body: []byte(fmt.Sprintf("resp: %s", m.Body))})
			m.Ack(false)
		}
	}()
	log.Println("waiting for msg....")
	select {}
}
