package mq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/oss"
)

// ConsumeMessage: consumes a message from RabbitMQ
func (r *RabbitMQ) ConsumeMessage() error {
	// Register a consumer
	msgs, err := r.channel.Consume(
		r.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			processMessage(msg.Body)
		}
	}()
	<-forever

	return nil
}

// processMessage: processes the message
func processMessage(message []byte) {
	var fileMsg FileTransferMessage
	err := json.Unmarshal(message, &fileMsg)
	if err != nil {
		log.Printf("failed to unmarshal the message: %v\n", err)
		return
	}

	// Upload the file to the OSS
	if err := oss.UploadFile(config.BucketName, fileMsg.ObjectKey, fileMsg.LocalFile); err != nil {
		log.Printf("failed to upload the file to the OSS: %v\n", err)
		return
	}
}
