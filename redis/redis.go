package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

// Client contains Redis client
var Client *redis.Client
var ctx = context.Background()

// Nil contains the nil value for Redis
var Nil = redis.Nil

// Connect function creates connection to the Redis server
func Connect() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")

	// create a client instance
	Client = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	// ping the server
	_, pingError := Client.Ping(ctx).Result()
	if pingError != nil {
		return pingError
	}

	return nil
}
