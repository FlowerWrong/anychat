package db

import (
	"sync"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/spf13/viper"
)

var (
	mqClient *nats.Conn
	mqOnce   sync.Once
)

func initMQClient() (err error) {
	mqClient, err = nats.Connect(viper.Get("nats_url").(string), nats.Timeout(10*time.Second))
	if err != nil {
		return err
	}
	return nil
}

// MQClient return a mq client
func MQClient() *nats.Conn {
	if mqClient == nil {
		mqOnce.Do(func() {
			err := initMQClient()
			if err != nil {
				panic(err)
			}
		})
	}
	return mqClient
}
