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
	const concurrency = 5
	var waitGroup sync.WaitGroup

	moviesChannel := make(chan struct {
		movieID int
		movie   domain.Movie
	}, len(movies))

	go func() {
		for movieID, movie := range movies {
			moviesChannel <- struct {
				movieID int
				movie   domain.Movie
			}{movieID, movie}
		}
		close(moviesChannel)
	}()

	for i := 0; i < concurrency; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			pipeline := redisClient.Pipeline()

			for result := range moviesChannel {
				movieAsJSON, err := json.Marshal(result.movie)
				if err != nil {
					log.Printf("Error marshaling movie into JSON: %v", err)
					continue
				}
				pipeline.JSONSet(ctx, KeyPrefix+strconv.Itoa(result.movieID+1), "$", string(movieAsJSON))
			}

			_, err := pipeline.Exec(ctx)
			if err != nil {
				log.Printf("Error writing JSON documents into Redis: %v", err)
			}
		}()
	}

	waitGroup.Wait()
	fmt.Printf("🟥 Movies stored on Redis: %d \n", len(movies))
}
