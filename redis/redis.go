package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Client *redis.Client
var ctx = context.Background()
var Nil = redis.Nil

// Connect to the Redis server
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
