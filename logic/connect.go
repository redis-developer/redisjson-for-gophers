package logic

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var redisConnectionURL = func() string {
	if redisConnURL := os.Getenv("REDIS_CONNECTION_URL"); redisConnURL != "" {
		return redisConnURL
	}
	return "redis://localhost:6379"
}()

func ConnectWithRedis(ctx context.Context) (*redis.Client, error) {
	connOpts, err := redis.ParseURL(redisConnectionURL)
	if err != nil {
		log.Printf("Error parsing the Redis URL: %v", err)
		return nil, err
	}
	connOpts.UnstableResp3 = true // Required for Search
	redisClient := redis.NewClient(connOpts)

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Ping request failed: %v", err)
		return nil, err
	}

	return redisClient, nil
}
