package mq

import "github.com/Shopify/sarama"

// 将信息放入 kafka

var kafkaProducer sarama.AsyncProducer

var brokerList = []string{"127.0.0.1:9010", "127.0.0.1:9091"}

// TODO 使用 kafka 时候需要重新添加 kafka
func init() {
	//var err error
	//config := sarama.NewConfig()
	//config.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Retry.Max = 10
	//config.Producer.Return.Successes = true
	//if kafkaProducer, err = sarama.NewAsyncProducer(brokerList, config); err != nil {
	//	panic(err)
	//}
}

func PublishKafka(topic, key, msg string) {
	kafkaProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		// 通过 key 来将数据映射到某个特定的 partition
		Key: sarama.StringEncoder(msg),
		Value: sarama.StringEncoder(msg),
	}
}
