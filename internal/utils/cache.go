package utils

import (
	"TeamTrackerBE/internal/config"
	"time"
)

func SetCache(key string, value string, ttl time.Duration) error {
	return config.RedisClient.Set(config.Ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return config.RedisClient.Get(config.Ctx, key).Result()
}

func DeleteCache(key string) error {
	return config.RedisClient.Del(config.Ctx, key).Err()
}
