package config

import (
	"log"
	"os"
	"time"

	"github.com/bladewaltz9/file-store-server/utils"
)

var (
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
)

const (
	BucketDir     = "file-store/"
	URLExpireTime = time.Hour * 24 // 24 hours
)

func init() {
	// Load the environment variables
	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Failed to load the .env file: %v", err)
	}

	// OSS
	Endpoint = os.Getenv("OSS_ENDPOINT")
	AccessKeyID = os.Getenv("OSS_ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("OSS_ACCESS_KEY_SECRET")
	BucketName = os.Getenv("OSS_BUCKET_NAME")
}
