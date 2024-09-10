package config

import (
	"log"
	"os"
	"strconv"

	"github.com/bladewaltz9/file-store-server/utils"
)

var (
	RabbitMQHost     string
	RabbitMQPort     int
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQVHost    string
	RabbitMQURL      string

	TransExchangeName  string
	TransOSSQueueName  string
	TransOSSRoutingKey string
)

func init() {
	// Load the environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load the .env file: %v", err)
	}

	// RabbitMQ
	RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	RabbitMQPort, _ = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	RabbitMQUser = os.Getenv("RABBITMQ_USER")
	RabbitMQPassword = os.Getenv("RABBITMQ_PASSWORD")
	RabbitMQVHost = os.Getenv("RABBITMQ_VHOST")

	RabbitMQURL = "amqp://" + RabbitMQUser + ":" + RabbitMQPassword + "@" + RabbitMQHost + ":" + strconv.Itoa(RabbitMQPort) + "/" + RabbitMQVHost

	TransExchangeName = os.Getenv("TRANS_EXCHANGE_NAME")
	TransOSSQueueName = os.Getenv("TRANS_OSS_QUEUE_NAME")
	TransOSSRoutingKey = os.Getenv("TRANS_OSS_ROUTING_KEY")
}
