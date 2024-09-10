package config

import (
	"log"
	"os"
	"time"

	"github.com/bladewaltz9/file-store-server/utils"
)

var (
	JWTSecretKey string
)

const (
	JWTExpirationTime = time.Hour * 1 // 1 hour
)

func init() {
	// Load the environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load the .env file: %v", err)
	}

	// JWT
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
}
