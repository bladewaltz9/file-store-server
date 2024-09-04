package redis

import (
	"context"
	"fmt"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/redis/go-redis/v9"
)

var (
	// redisClient: redis client
	redisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis: initialize the redis client
func initRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the redis: %v", err.Error()))
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func init() {
	// initialize the redis client
	addr := config.RedisHost + ":" + fmt.Sprintf("%d", config.RedisPort)
	password := config.RedisPassword
	db := config.RedisDB

	initRedis(addr, password, db)
}
