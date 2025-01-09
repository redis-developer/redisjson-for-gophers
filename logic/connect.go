package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"redisjson4gophers/domain"
)

var redisConnectionURL = func() string {
	if redisConnURL := os.Getenv("REDIS_CONNECTION_URL"); redisConnURL != "" {
		return redisConnURL
	}
	return "redis://localhost:6379"
}()

func ConnectWithRedis(ctx context.Context) context.Context {
	connOpts, err := redis.ParseURL(redisConnectionURL)
	if err != nil {
		panic(fmt.Errorf("failed to parse Redis URL: %w", err))
	}
	connOpts.UnstableResp3 = true // Required for Search
	redisClient := redis.NewClient(connOpts)

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("error connecting with Redis: %w", err))
	}

	return context.WithValue(ctx, domain.ClientKey, redisClient)
}
