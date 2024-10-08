package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

// maxSubDirCount: the maximum subdirectory count
const maxSubDirCount = 10

var envLoaded = false

// LoadEnv: load the environment variables
func LoadEnv() error {
	if envLoaded {
		return nil
	}
	path := "config/.env"
	for i := 0; i < maxSubDirCount; i++ {
		err := godotenv.Load(path)
		if err == nil {
			return nil
		}
		path = "../" + path
	}
	return fmt.Errorf("failed to load the .env file")
}
