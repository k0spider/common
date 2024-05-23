package mq

import (
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	*RabbitMQInfoConfig
	delayQueueName     string
	delayQueueRouteKey string
}

type Consumer func(amqp.Delivery)

func newRabbitMQ(ch *amqp.Channel, config *RabbitMQInfoConfig) (*RabbitMQ, error) {
	rabbitMQ := &RabbitMQ{RabbitMQInfoConfig: config, ch: ch}
	err := rabbitMQ.ch.ExchangeDeclare(config.ExchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = rabbitMQ.claimQueue(rabbitMQ.QueueName, rabbitMQ.RouteKey, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return rabbitMQ, nil
}

func (mq *RabbitMQ) createDelayQueue() error {
	if mq.delayQueueRouteKey != "" {
		return nil
	}
	mq.delayQueueName = mq.QueueName + "-delay"
	mq.delayQueueRouteKey = mq.RouteKey + "-delay"
	err := mq.claimQueue(mq.delayQueueName, mq.delayQueueRouteKey, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    mq.ExchangeName,
		"x-dead-letter-routing-key": mq.RouteKey,
	})
	return err
}

func (mq *RabbitMQ) claimQueue(queueName, routeKey string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) error {
	_, err := mq.ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, args)
	if err != nil {
		return err
	}
	err = mq.ch.QueueBind(queueName, routeKey, mq.ExchangeName, false, nil)
	if err != nil {
		return err
	}
	return err
}

// 发送延迟消息
func (mq *RabbitMQ) sendDelayMessage(publishing *amqp.Publishing) error {
	if err := mq.createDelayQueue(); err != nil {
		return err
	}
	return mq.ch.Publish(mq.ExchangeName, mq.delayQueueRouteKey, false, false, *publishing)
}

// 发送消息
func (mq *RabbitMQ) sendMessage(publishing *amqp.Publishing) error {
	return mq.ch.Publish(mq.ExchangeName, mq.RouteKey, false, false, *publishing)
}

// 读取消息
func (mq *RabbitMQ) consumers() (<-chan amqp.Delivery, error) {
	err := mq.claimQueue(mq.QueueName, mq.RouteKey, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = mq.ch.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}
	return mq.ch.Consume(mq.QueueName, "", false, false, false, false, nil)
}
