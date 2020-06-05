package app

import (
	"fmt"
	config "github.com/stetsd/monk-conf"
)

type Sender struct {
	config config.Config
}

func NewApp(config config.Config) *Sender {
	return &Sender{config: config}
}

func (sender *Sender) Start() {
	fmt.Println("START")
}

func (sender *Sender) Stop() {
	fmt.Println("STOP")
}
