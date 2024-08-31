package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	DBHost = os.Getenv("MYSQL_HOST")
	DBPort, _ = strconv.Atoi(os.Getenv("MYSQL_PORT"))
	DBUser = os.Getenv("MYSQL_USER")
	DBPassword = os.Getenv("MYSQL_PASSWORD")
	DBName = os.Getenv("MYSQL_DATABASE")

	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
}
