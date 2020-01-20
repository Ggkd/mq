package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll	// 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner	// 新选出一个partition
	config.Producer.Return.Successes = true	// 成功交付的消息将在success channel返回

	//构造一个消息
	message := &sarama.ProducerMessage{}
	message.Topic = "api_test"
	message.Value = sarama.StringEncoder("this is test value")
	//连接kafka
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()
	//发送消息
	for {
		partition, offset, err := producer.SendMessage(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("partition : %v,  offset : %v\n", partition, offset)
		time.Sleep(time.Millisecond * 500)
	}
}