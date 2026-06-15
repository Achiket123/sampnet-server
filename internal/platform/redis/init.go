package redis

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if addr == "" {
		addr = "localhost:6379"
	}

	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			log.Fatalf("Invalid REDIS_DB configuration: %v", err)
		}
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping check
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Printf("Connected to Redis successfully at %s", addr)
}
