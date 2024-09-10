package mq_test

import (
	"testing"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/mq"
)

// TestRabbitMQ: tests the RabbitMQ
func TestRabbitMQ(t *testing.T) {
	fileMsg := &mq.FileTransferMessage{
		FileID:    1,
		LocalFile: "/home/bladewaltz/workspace/go/file-store-server/main.go",
		ObjectKey: config.BucketDir + "main.go",
	}

	rabbitMQ, err := mq.NewRabbitMQ(config.TransExchangeName, config.TransOSSQueueName, config.TransOSSRoutingKey, config.RabbitMQURL)
	if err != nil {
		t.Errorf("failed to create a new RabbitMQ instance: %v", err)
	}

	// Publish a message
	if err := rabbitMQ.PublishMessage(fileMsg); err != nil {
		t.Errorf("failed to publish a message: %v", err)
	}

	// Consume a message
	if err := rabbitMQ.ConsumeMessage(); err != nil {
		t.Errorf("failed to consume a message: %v", err)
	}
}
