package utils

import (
	"TeamTrackerBE/internal/config"
	"time"
)

// SetCache stores a value in Redis using the given key and TTL (time-to-live).
func SetCache(key string, value string, ttl time.Duration) error {
	return config.RedisClient.Set(config.Ctx, key, value, ttl).Err()
}

// GetCache retrieves a value from Redis by key.
func GetCache(key string) (string, error) {
	return config.RedisClient.Get(config.Ctx, key).Result()
}

// DeleteCache removes a value from Redis by key.
func DeleteCache(key string) error {
	return config.RedisClient.Del(config.Ctx, key).Err()
}

// IsCacheAvailable reports whether a Redis client is configured.
func IsCacheAvailable() bool {
	return config.RedisClient != nil
}

// ListRPush appends one or more values to a Redis list (RPUSH).
func ListRPush(key string, values ...string) error {
	if len(values) == 0 {
		return nil
	}
	any := make([]interface{}, 0, len(values))
	for _, v := range values {
		any = append(any, v)
	}
	return config.RedisClient.RPush(config.Ctx, key, any...).Err()
}

// ListLRange returns list items between start and stop (inclusive).
func ListLRange(key string, start, stop int64) ([]string, error) {
	return config.RedisClient.LRange(config.Ctx, key, start, stop).Result()
}

// Expire sets TTL for a key.
func Expire(key string, ttl time.Duration) error {
	return config.RedisClient.Expire(config.Ctx, key, ttl).Err()
}

// TTL returns time-to-live for a key.
func TTL(key string) (time.Duration, error) {
	return config.RedisClient.TTL(config.Ctx, key).Result()
}
