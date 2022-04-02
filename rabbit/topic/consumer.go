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
	// 3.声明交换机
	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic", // 交换机类型
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("declear exchange failed: %s\n", err)
	}
	// 4.声明队列
	queue, err := ch.QueueDeclare(
		"",
		false, // 非持久化队列
		false,
		true, //独占队列（当前声明的队列的连接断开时删除）
		false,
		nil)
	if err != nil {
		log.Printf("declare queue failed: %s\n", err)
	}

	// 5.绑定队列
	ch.QueueBind(
		queue.Name,   //队列名称
		"*.info.#",   // routing-key
		"logs_topic", // 交换机
		false,
		nil,
	)

	// 5.绑定队列
	ch.QueueBind(
		queue.Name,   //队列名称
		"#.error",    // routing-key
		"logs_topic", // 交换机
		false,
		nil,
	)

	// 6.获取消息
	deliver, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Printf("register consumer failed: %s\n", err)
	}
	go func() {
		for msg := range deliver {
			log.Printf("receive msg: %s, from exchange: %s , routingKey:%s,  queue:%s, \n", msg.Body, msg.Exchange, msg.RoutingKey, queue.Name)
			log.Println("Done!")
		}
	}()
	log.Println("waiting for msg....")
	select {}
}
