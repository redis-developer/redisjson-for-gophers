package main

import (
	"context"
	"redisjson4gophers/logic"
)

func main() {

	ctx := context.Background()

	ctx = logic.LoadMoviesFromFile(ctx)
	ctx = logic.ConnectWithRedis(ctx)
	logic.IndexMoviesAsDocuments(ctx)
	logic.LookupMovieTitleByMovieKey(ctx)
	logic.SearchBestMatrixMovies(ctx)
	logic.MovieCountPerGenreAgg(ctx)
}
