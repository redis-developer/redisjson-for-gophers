package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"redisjson4gophers/domain"
	"strconv"
)

func IndexMoviesAsDocuments(ctx context.Context, redisClient *redis.Client, movies []domain.Movie) {
	pipeline := redisClient.Pipeline()
	for movieID, movie := range movies {
		movieAsJSON, err := json.Marshal(movie)
		if err != nil {
			log.Printf("Error marshaling movie into JSON: %v", err)
		}
		pipeline.JSONSet(ctx, keyPrefix+strconv.Itoa(movieID+1), "$", string(movieAsJSON))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Printf("Error writing JSON documents into Redis: %v", err)
	}

	fmt.Printf("🟥 Movies stored on Redis: %d \n", len(movies))
}
