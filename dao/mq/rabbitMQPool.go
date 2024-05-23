package mq

import (
	"context"
	"github.com/k0spider/common/log"
	"github.com/kirinlabs/rabbitgo"
	"github.com/streadway/amqp"
	"time"
)

var rabbitMQPool *rabbitgo.RabbitPool

type RabbitMQInfoConfig struct {
	ExchangeName string `yaml:"exchangeName"`
	QueueName    string `yaml:"queueName"`
	RouteKey     string `yaml:"routeKey"`
}

type RabbitMQPoolConfig struct {
	Dsn   string `yaml:"dsn"`
	Debug bool   `yaml:"debug"`
}

func InitPool(config *RabbitMQPoolConfig) {
	rabbitMQPool = rabbitgo.New(config.Dsn, rabbitgo.Config{
		ConnectionMax: 5,
		ChannelMax:    10,
		ChannelActive: 20,
		ChannelIdle:   10,
	})
	if !config.Debug {
		rabbitMQPool.SetLevel(3)
	}
}

func getChannel(config *RabbitMQInfoConfig) (*rabbitgo.Channel, *RabbitMQ, error) {
	ch, err := rabbitMQPool.Get()
	if err != nil {
		return nil, nil, err
	}
	rabbitMQ, err := newRabbitMQ(ch.Ch, config)
	if err != nil {
		return nil, nil, err
	}
	return ch, rabbitMQ, nil
}

func Push(ctx context.Context, config *RabbitMQInfoConfig, publishing *amqp.Publishing) error {
	ch, rabbitMQ, err := getChannel(config)
	if err != nil {
		log.WithContext(ctx).Errorf("GetCh err:%v", err)
		return err
	}
	defer rabbitMQPool.Push(ch)
	if publishing.Expiration != "" {
		err = rabbitMQ.sendDelayMessage(publishing)
	} else {
		err = rabbitMQ.sendMessage(publishing)
	}
	if err != nil {
		log.WithContext(ctx).Errorf("SendMessage err:%v", err)
		return err
	}
	return nil
}

// Tell me your handling directly
func Consumers(ctx context.Context, config *RabbitMQInfoConfig, fn Consumer) error {
	_, _, err := getChannel(config)
	if err == nil {
		go listen(*config, fn)
	}
	return err
}

// How do you deal with it? It's up to you
func GetMsgChannel(config *RabbitMQInfoConfig) (<-chan amqp.Delivery, error) {
	_, rabbitMQ, err := getChannel(config)
	if err != nil {
		return nil, err
	}
	return rabbitMQ.consumers()
}

func listen(config RabbitMQInfoConfig, fn Consumer) {
	ctx := context.Background()
	defer func() {
		time.Sleep(time.Second * 3)
		if panicMsg := recover(); panicMsg != nil {
			log.WithContext(ctx).Errorf("MQ Abnormal operation:%v", panicMsg)
		}
		go listen(config, fn)
	}()
	channelMsg, err := GetMsgChannel(&config)
	if err != nil {
		log.WithContext(ctx).Errorf("MQ GetCh err:%v", err)
		return
	}
	for s := range channelMsg {
		fn(s)
	}
}
