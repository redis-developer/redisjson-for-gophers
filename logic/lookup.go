package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"redisjson4gophers/domain"
	"strconv"
)

func LookupMovieTitleByMovieKey(ctx context.Context) {
	movies := ctx.Value(domain.MoviesKey).([]domain.Movie)
	redisClient := ctx.Value(domain.ClientKey).(*redis.Client)

	movieKey := "movie:" + strconv.Itoa(rand.Intn(len(movies)))
	movieTitle, err := redisClient.JSONGet(ctx, movieKey, "$.title").Expanded()
	if err != nil {
		log.Printf("Error getting the movie title: %v", err)
	}

	fmt.Printf("🟥 Movie with the key '%s': %s \n", movieKey, movieTitle.([]interface{})[0])
}
