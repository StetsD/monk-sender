package app

import (
	"fmt"
	config "github.com/stetsd/monk-conf"
	"github.com/stetsd/monk-sender/internal/app/contracts"
	"github.com/stetsd/monk-sender/internal/infrastructure"
)

type Sender struct {
	config      config.Config
	queueClient contracts.QueueClient
}

func NewApp(config config.Config) *Sender {
	return &Sender{config: config}
}

func (sender *Sender) Start() {
	fmt.Println("START")
	qC := infrastructure.NewKafkaClient()
	sender.queueClient = qC

	sender.queueClient.Init()
}

func (sender *Sender) Stop() {
	fmt.Println("STOP")
}
