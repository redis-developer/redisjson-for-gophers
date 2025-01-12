package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"strconv"
)

func LookupMovieTitleByMovieKey(ctx context.Context, redisClient *redis.Client, arrLength int) {
	movieKey := "movie:" + strconv.Itoa(rand.Intn(arrLength))
	movieTitle, err := redisClient.JSONGet(ctx, movieKey, "$.title").Result()
	if err != nil {
		log.Printf("Error getting the movie title: %v", err)
	}

	fmt.Printf("🟥 Movie with the key '%s': %s \n", movieKey, movieTitle)
}
