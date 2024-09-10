package mq

import (
	"fmt"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/streadway/amqp"
)

// RabbitMQ: represents the RabbitMQ connection
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	Exchange string
	Queue    string
	Key      string
	url      string
}

var rabbitMQ *RabbitMQ

// NewRabbitMQ: creates a new RabbitMQ instance
func NewRabbitMQ(exchange, queue, key, url string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		Exchange: exchange,
		Queue:    queue,
		Key:      key,
		url:      url,
	}

	if err := rmq.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	return rmq, nil
}

func init() {
	var err error
	rabbitMQ, err = NewRabbitMQ(config.TransExchangeName, config.TransOSSQueueName, config.TransOSSRoutingKey, config.RabbitMQURL)
	if err != nil {
		panic(fmt.Sprintf("failed to create a new RabbitMQ instance: %v", err))
	}
}

func GetRabbitMQ() *RabbitMQ {
	return rabbitMQ
}

// connect: connects to RabbitMQ
func (r *RabbitMQ) connect() error {
	var err error
	r.conn, err = amqp.Dial(r.url)
	if err != nil {
		return err
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return err
	}
	return nil
}
