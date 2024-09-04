package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/bladewaltz9/file-store-server/redis"
)

var ctx = context.Background()

func TestRedisSetGet(t *testing.T) {
	redisClient := redis.GetRedisClient()

	// test set
	err := redisClient.Set(ctx, "test_key", "test_value", 0).Err()
	if err != nil {
		t.Errorf("Failed to set the key: %v", err.Error())
	}

	// test get
	value, err := redisClient.Get(ctx, "test_key").Result()
	if err != nil {
		t.Errorf("Failed to get the key: %v", err.Error())
	}
	if value != "test_value" {
		t.Errorf("The value is not correct: %v", value)
	}

	// clean up
	if err := redisClient.Del(ctx, "test_key").Err(); err != nil {
		t.Errorf("Failed to delete the key: %v", err.Error())
	}
}

func TestRedisExpiration(t *testing.T) {
	redisClient := redis.GetRedisClient()

	// test set with expiration
	err := redisClient.Set(ctx, "test_key", "test_value", 1*time.Second).Err()
	if err != nil {
		t.Errorf("Failed to set the key with expiration: %v", err.Error())
	}

	// test get before expiration
	value, err := redisClient.Get(ctx, "test_key").Result()
	if err != nil {
		t.Errorf("Failed to get the key before expiration: %v", err.Error())
	}
	if value != "test_value" {
		t.Errorf("The value before expiration is not correct: %v", value)
	}

	// test get after expiration
	time.Sleep(2 * time.Second)
	value, err = redisClient.Get(ctx, "test_key").Result()
	if err == nil {
		t.Errorf("The key should be expired, but got the value: %v", value)
	}
}
