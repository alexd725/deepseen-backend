package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Client RedisClient
var ctx = context.Background()

// Connect to the Redis server
func Connect() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")

	// create a client instance
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	// ping the server
	_, pingError := client.Ping(ctx).Result()
	if pingError != nil {
		return pingError
	}

	Client = client
	return nil
}
