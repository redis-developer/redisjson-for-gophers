package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"redisjson4gophers/logic"
)

func main() {
	ctx := context.Background()

	redisClient, err := logic.ConnectWithRedis(ctx)
	if err != nil {
		log.Fatalf("Error connecting with Redis: %v", err)
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Fatalf("Error closing the Redis connection: %v", err)
		}
	}(redisClient)

	movies, err := logic.LoadMoviesFromFile("movies.json")
	if err != nil {
		log.Fatalf("Error loading movies from file: %v", err)
	}

	logic.IndexMoviesAsDocuments(ctx, redisClient, movies)
	logic.LookupMovieTitleByMovieKey(ctx, redisClient, len(movies))
	logic.SearchBestMatrixMovies(ctx, redisClient)
	logic.MovieCountPerGenreAgg(ctx, redisClient)
}
