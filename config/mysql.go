package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bladewaltz9/file-store-server/utils"
)

var (
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
)

const (
	// database connection pool settings
	DBMaxConn         = 100
	DBMaxIdleConn     = 30
	DBConnMaxLifetime = time.Hour
)

func init() {
	// Load the environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load the .env file: %v", err)
	}

	// MySQL
	DBHost = os.Getenv("MYSQL_HOST")
	DBPort, _ = strconv.Atoi(os.Getenv("MYSQL_PORT"))
	DBUser = os.Getenv("MYSQL_USER")
	DBPassword = os.Getenv("MYSQL_PASSWORD")
	DBName = os.Getenv("MYSQL_DATABASE")
}
