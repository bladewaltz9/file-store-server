package redis

import (
	"context"
	"fmt"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/redis/go-redis/v9"
)

var (
	// rdb: redis client
	rdb *redis.Client
	ctx = context.Background()
)

// InitRedis: initialize the redis client
func initRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the redis: %v", err.Error()))
	}
}

func GetRedisClient() *redis.Client {
	return rdb
}

func init() {
	// initialize the redis client
	addr := config.RedisHost + ":" + fmt.Sprintf("%d", config.RedisPort)
	password := config.RedisPassword
	db := config.RedisDB

	initRedis(addr, password, db)
}
