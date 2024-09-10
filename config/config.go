package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bladewaltz9/file-store-server/utils"
)

// Directory path
const (
	FileStoreDir = "/home/bladewaltz/data/files/"
	FileChunkDir = "/home/bladewaltz/data/chunks/"
)

// Upload size limit
const MaxUploadSize = 32 << 20 // 32MB

// MySQL config
var (
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// database connection pool settings
	DBMaxConn         = 100
	DBMaxIdleConn     = 30
	DBConnMaxLifetime = time.Hour
)

// Redis config
var (
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
)

// JWT config
var (
	JWTSecretKey      string
	JWTExpirationTime = time.Hour * 1 // 1 hour
)

// SSL cert and key
const (
	CertFile = "/etc/apache2/ssl/bladewaltz.cn.crt"
	KeyFile  = "/etc/apache2/ssl/bladewaltz.cn.key"
)

// OSS config
var (
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	BucketDir       = "file-store/"
	URLExpireTime   = time.Hour * 24 // 24 hours
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

	// Redis
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort, _ = strconv.Atoi(os.Getenv("REDIS_PORT"))
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

	// JWT
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	// OSS
	Endpoint = os.Getenv("OSS_ENDPOINT")
	AccessKeyID = os.Getenv("OSS_ACCESS_KEY_ID")
	AccessKeySecret = os.Getenv("OSS_ACCESS_KEY_SECRET")
	BucketName = os.Getenv("OSS_BUCKET_NAME")
}
