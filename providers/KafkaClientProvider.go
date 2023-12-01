package providers

import (
	"encoding/json"
	"fmt"
	"jericho-go/types"
	"jericho-go/wrongs"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

var wg sync.WaitGroup

func KafkaClientHandler(addr, topic string) {
	var err error

	consumer, err := sarama.NewConsumer([]string{addr}, nil)

	if err != nil {
		fmt.Println("[kafka-client] 创建消费者失败：", err)
		return
	}
	fmt.Println("[kafka-client] 链接成功")
	defer func(consumer sarama.Consumer) {
		err := consumer.Close()
		if err != nil {
			wrongs.ThrowForbidden("[kafka-client] 关闭失败")
		}
	}(consumer)
	partitions, err := consumer.Partitions(topic)

	if err != nil {
		fmt.Println("[kafka-client] 获取消息失败：", err)
		return
	}

	for _, p := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			fmt.Println("[kafka-client] 消费数据失败：", err)
			continue
		}
		wg.Add(1)
		go func() {
			for m := range partitionConsumer.Messages() {
				fmt.Printf("[kafka-client] 接收消息 partitions: %d, key: %s, text: %s, offset: %d\n", m.Partition, string(m.Key), string(m.Value), m.Offset)

				business := &types.StdBusiness{}
				if err = json.Unmarshal(m.Value, business); err != nil {
					log.Printf("[kafka-client] [解析业务失败] %v", err)
				}

				switch business.BusinessType {
				case "ping":
					log.Printf("[kafka-client] [ping]")
					// RabbitMqSendMessage(tools.NewCorrectWithBusiness("pong", "pong").Datum(map[string]any{"time": time.Now().Unix()}).ToJsonStr())
				case "message":
					log.Printf("[kafka-client] [message]")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
