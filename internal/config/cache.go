package config

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisClient *redis.Client

func SetupCache() {
	env := Env.Cache

	db, err := strconv.Atoi(env.DB)
	if err != nil {
		log.Fatalf("Invalid Redis DB index: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: env.Host+":"+env.Port,
		Password: env.Password,
		DB: db,
	})

	if _, err := RedisClient.Ping(Ctx).Result(); err != nil {
		log.Fatalf(" Failed to connect Redis: %v", err)
	}

	log.Println("Redis connected")
}
