package contracts

import "github.com/Shopify/sarama"

type TransportClient interface {
	InitConsumer(topic string) (sarama.PartitionConsumer, error)
	InitProducer() (sarama.AsyncProducer, error)
}
