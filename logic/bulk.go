package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"redisjson4gophers/domain"
	"strconv"
	"sync"
)

func IndexMoviesAsDocuments(ctx context.Context, redisClient *redis.Client, movies []domain.Movie) {
	var waitGroup sync.WaitGroup
	moviesChannel := make(chan struct {
		id    int
		movie domain.Movie
	}, len(movies))

	for movieID, movie := range movies {
		waitGroup.Add(1)
		go func(movieID int, movie domain.Movie) {
			defer waitGroup.Done()
			movie.PlotEmbedding = CreateEmbedding(ctx, movie.Plot)
			moviesChannel <- struct {
				id    int
				movie domain.Movie
			}{movieID, movie}
		}(movieID, movie)
	}

	go func() {
		waitGroup.Wait()
		close(moviesChannel)
	}()

	pipeline := redisClient.Pipeline()
	for result := range moviesChannel {
		movieAsJSON, err := json.Marshal(result.movie)
		if err != nil {
			log.Printf("Error marshaling movie into JSON: %v", err)
			continue
		}
		pipeline.JSONSet(ctx, KeyPrefix+strconv.Itoa(result.id+1), "$", string(movieAsJSON))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Printf("Error writing JSON documents into Redis: %v", err)
	}

	fmt.Printf("🟥 Movies stored on Redis: %d \n", len(movies))
}
