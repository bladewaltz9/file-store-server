package mq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// PublishMessage: publishes a message to RabbitMQ
func (r *RabbitMQ) PublishMessage(fileMsg *FileTransferMessage) error {
	// Serialize the message to JSON
	body, err := json.Marshal(fileMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal the message: %v", err)
	}

	// Publish the message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish the message: %v", err)
	}
	return nil
}
