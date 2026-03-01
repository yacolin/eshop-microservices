package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// NewRedis 创建 Redis 客户端
func NewRedis(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	log.Println("redis connected")
	return client, nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// GetCtx returns the context for Redis operations
func GetCtx() context.Context {
	return ctx
}

// GetDuration returns time.Duration from seconds
func GetDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// GetCache gets value from Redis cache by key
func GetCache(key string) ([]byte, error) {
	return client.Get(ctx, key).Bytes()
}

// SetCache sets value to Redis cache with expiration
func SetCache(key string, value []byte, expiration int) error {
	return client.SetEx(ctx, key, value, GetDuration(expiration)).Err()
}

// DeleteCache deletes value from Redis cache by key
func DeleteCache(key string) error {
	return client.Del(ctx, key).Err()
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	// In a real implementation, you would use os.Getenv(key)
	// For now, we'll just return the default value
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	// In a real implementation, you would parse os.Getenv(key)
	// For now, we'll just return the default value
	return defaultValue
}
