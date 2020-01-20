package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)

func main() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	partitionList, err := consumer.Partitions("log_path")		//根据topic找到所有的分区
	if err != nil {
		fmt.Println(err)
		return
	}
	//遍历所有的分区
	wg := sync.WaitGroup{}
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition("log_path", partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err)
		}
		defer pc.Close()
		//异步的消费数据
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages(){
				fmt.Printf("topic:%s,  partiotion:%v, offset:%v, value:%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}