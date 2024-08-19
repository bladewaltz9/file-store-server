package config

import "time"

const (
	FileStorePath = "/tmp/file-store-server/"
	MaxUploadSize = 32 << 20 // 32MB

	// MySQL config
	DBHost     = "127.0.0.1"
	DBPort     = 3306
	DBUser     = "root"
	DBPassword = "Lollzp1999!"
	DBName     = "file_server"
	DBMaxConn  = 1000

	// JWT config
	JWTSecretKey      = "BxHpBYB3rey1bOidVbcCiHa389t5edWkW7yo1vPLXxc="
	JWTExpirationTime = time.Hour * 1 // 1 hour
)
