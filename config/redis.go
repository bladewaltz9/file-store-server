package config

import (
	"log"
	"os"
	"strconv"

	"github.com/bladewaltz9/file-store-server/utils"
)

var (
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
)

func init() {
	// Load the environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load the .env file: %v", err)
	}

	// Redis
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
}
