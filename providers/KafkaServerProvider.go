package providers

import (
	"fmt"
	"jericho-gin/wrongs"

	"github.com/IBM/sarama"
)

var KafkaSyncProducer sarama.SyncProducer

func KafkaServerHandler(addr, topic string) {
	launchKafkaServer(addr)
	// 定义一个生产消息，包括Topic、消息内容
	// KafkaSendMessage(topic, "a", "a4")
}

func launchKafkaServer(addr string) {
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true

	// 新建一个同步生产者
	KafkaSyncProducer, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		fmt.Println("[kafka-server] 创建生产者失败：", err)
		return
	}
}

func CloseKafkaServer() {
	if KafkaSyncProducer != nil {
		if err := KafkaSyncProducer.Close(); err != nil {
			wrongs.ThrowForbidden("[kafka-server] 关闭失败")
		}
		fmt.Println("[kafka-server] 关闭成功")
	}
}

// KafkaSendMessage 发送kafka消息
func KafkaSendMessage(topic, key, content string) error {
	// 定义一个生产消息，包括Topic、消息内容
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}

	// 发送消息
	pid, offset, err := KafkaSyncProducer.SendMessage(msg)
	if err != nil {
		fmt.Println(fmt.Sprintf("[kafka-server] 发送消息失败： %s", err.Error()))
		return err
	} else {
		fmt.Println(fmt.Sprintf("[kafka-server] 发送消息成功：pid「%d」 offset「%d」", pid, offset))
		return nil
	}
}
