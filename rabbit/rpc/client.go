package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
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

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

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
	queue, err := ch.QueueDeclare(
		"",
		false, // 持久化队列
		false,
		true,
		false,
		nil)
	if err != nil {
		log.Printf("declare queue failed: %s\n", err)
	}

	// 消费消息
	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("receive register consumer failed: %s\n", err)
	}

	corrId := randomString(32)
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ReplyTo:       queue.Name,
			CorrelationId: corrId,
			Body:          []byte(body),
		})
	if err != nil {
		log.Printf("publish msg failed\": %s\n", err)
	}
	for m := range msgs {
		if corrId == m.CorrelationId {
			fmt.Printf("callback msg: %s\n", m.Body)
			break
		}
	}
}
