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

const (
	indexName = "json_movies_index"
	keyPrefix = "movie:"
)

func IndexMoviesAsDocuments(ctx context.Context) {

	movies := ctx.Value(domain.MoviesKey).([]domain.Movie)
	redisClient := ctx.Value(domain.ClientKey).(*redis.Client)

	/************************************ Index Creation ************************************/

	redisClient.FTDropIndexWithArgs(ctx, indexName, &redis.FTDropIndexOptions{DeleteDocs: true})

	titleField := &redis.FieldSchema{FieldName: "$.title", FieldType: redis.SearchFieldTypeText, As: "title"}
	yearField := &redis.FieldSchema{FieldName: "$.year", FieldType: redis.SearchFieldTypeNumeric, As: "year"}
	plotField := &redis.FieldSchema{FieldName: "$.plot", FieldType: redis.SearchFieldTypeText, As: "plot"}
	runningTimeField := &redis.FieldSchema{FieldName: "$.runningTime", FieldType: redis.SearchFieldTypeNumeric, As: "runningTime"}
	releaseDateField := &redis.FieldSchema{FieldName: "$.releaseDate", FieldType: redis.SearchFieldTypeText, As: "releaseDate"}
	ratingField := &redis.FieldSchema{FieldName: "$.rating", FieldType: redis.SearchFieldTypeNumeric, As: "rating"}
	genresField := &redis.FieldSchema{FieldName: "$.genres.*", FieldType: redis.SearchFieldTypeTag, As: "genres"}
	actorsField := &redis.FieldSchema{FieldName: "$.actors.*", FieldType: redis.SearchFieldTypeTag, As: "actors"}
	directorsField := &redis.FieldSchema{FieldName: "$.directors.*", FieldType: redis.SearchFieldTypeTag, As: "directors"}

	redisClient.FTCreate(ctx, indexName,
		&redis.FTCreateOptions{OnJSON: true, Prefix: []interface{}{keyPrefix}},
		titleField, yearField, plotField, runningTimeField, releaseDateField,
		ratingField, genresField, actorsField, directorsField).Result()

	/******************************************************************************************/

	pipeline := redisClient.Pipeline()
	for movieID, movie := range movies {
		movieAsJSON, err := json.Marshal(movie)
		if err != nil {
			log.Printf("Error marshaling movie into JSON: %v", err)
		}
		pipeline.JSONSet(ctx, keyPrefix+strconv.Itoa(movieID), "$", string(movieAsJSON))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Printf("Error writing JSON documents into Redis: %v", err)
	}

	fmt.Printf("🟥 Movies indexed on Redis: %d \n", len(movies))
}
