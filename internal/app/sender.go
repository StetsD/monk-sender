package app

import (
	config "github.com/stetsd/monk-conf"
	"github.com/stetsd/monk-sender/internal/app/contracts"
	"github.com/stetsd/monk-sender/internal/infrastructure"
	"github.com/stetsd/monk-sender/internal/infrastructure/logger"
	"os"
	"os/signal"
	"sync"
)

type Sender struct {
	config          config.Config
	transportClient contracts.TransportClient
	channel         SenderStrategy
}

func NewApp(config config.Config) *Sender {
	return &Sender{config: config}
}

func (sender *Sender) Start() {
	logger.Log.Info("Service monk-sender is running")
	sender.transportClient = infrastructure.NewKafkaClient(sender.config)

	consumer, err := sender.transportClient.InitConsumer("on_send")
	if err != nil {
		panic(err)
	}

	sender.channel = &StrategySmtp{}
	if err := sender.channel.Init(map[string]string{
		"username": os.Getenv("SMTP_USER"),
		"password": os.Getenv("SMTP_PASS"),
		"host":     os.Getenv("SMTP_HOST"),
		"port":     os.Getenv("SMTP_PORT"),
	}); err != nil {
		logger.Log.Fatal(err.Error())
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Log.Fatal(err.Error())
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
	outer:
		for {
			select {
			case msg := <-consumer.Messages():
				received, err := onSendUnmarshaling(&msg.Value)
				if err != nil {
					logger.Log.Error(err.Error())
				}
				sender.channel.Send(received)
			case transportErr := <-consumer.Errors():
				logger.Log.Error(transportErr.Error())
			case <-signals:
				break outer
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
