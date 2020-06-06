package infrastructure

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
)

type KafkaClient struct {
	config sarama.Config
}

func NewKafkaClient() *KafkaClient {
	return &KafkaClient{}
}

func (kc *KafkaClient) Init() {

}

func (kc *KafkaClient) InitListener() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.SASL.Password = "12345"
	config.Net.SASL.User = "boris"
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("on_send", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
outer:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break outer
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}

func (kc *KafkaClient) InitProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.SASL.Password = "12345"
	config.Net.SASL.User = "boris"
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)

	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			panic(err)
			errors++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

outer:
	for {
		message := &sarama.ProducerMessage{Topic: "on_send", Value: sarama.StringEncoder("hello mazafaka")}
		select {
		case producer.Input() <- message:
			enqueued++
		case <-signals:
			producer.AsyncClose()
			break outer
		}
	}

	wg.Wait()
}

func (kc *KafkaClient) Send() {

}
