package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	IndexName = "movies_index"
	KeyPrefix = "movie:"
)

func CreateMoviesIndexOnRedis(ctx context.Context, redisClient *redis.Client) {
	redisClient.FTDropIndexWithArgs(ctx, IndexName, &redis.FTDropIndexOptions{DeleteDocs: true})

	titleField := &redis.FieldSchema{FieldName: "$.title", FieldType: redis.SearchFieldTypeText, As: "title"}
	yearField := &redis.FieldSchema{FieldName: "$.year", FieldType: redis.SearchFieldTypeNumeric, As: "year"}
	plotField := &redis.FieldSchema{FieldName: "$.plot", FieldType: redis.SearchFieldTypeText, As: "plot"}
	plotEmbeddingsField := &redis.FieldSchema{FieldName: "$.plotEmbedding", FieldType: redis.SearchFieldTypeVector, As: "plotEmbedding",
		VectorArgs: &redis.FTVectorArgs{
			HNSWOptions: &redis.FTHNSWOptions{Type: "FLOAT64", Dim: 1536, DistanceMetric: "COSINE"},
		},
	}
	runningTimeField := &redis.FieldSchema{FieldName: "$.runningTime", FieldType: redis.SearchFieldTypeNumeric, As: "runningTime"}
	releaseDateField := &redis.FieldSchema{FieldName: "$.releaseDate", FieldType: redis.SearchFieldTypeText, As: "releaseDate"}
	ratingField := &redis.FieldSchema{FieldName: "$.rating", FieldType: redis.SearchFieldTypeNumeric, As: "rating"}
	genresField := &redis.FieldSchema{FieldName: "$.genres.*", FieldType: redis.SearchFieldTypeTag, As: "genres", Separator: ","}
	actorsField := &redis.FieldSchema{FieldName: "$.actors.*", FieldType: redis.SearchFieldTypeTag, As: "actors", Separator: ","}
	directorsField := &redis.FieldSchema{FieldName: "$.directors.*", FieldType: redis.SearchFieldTypeTag, As: "directors", Separator: ","}

	_, err := redisClient.FTCreate(ctx, IndexName,
		&redis.FTCreateOptions{OnJSON: true, Prefix: []interface{}{KeyPrefix}},
		titleField, yearField, plotField, plotEmbeddingsField, runningTimeField,
		releaseDateField, ratingField, genresField, actorsField, directorsField).Result()
	if err != nil {
		log.Printf("Error creating the index: %v", err)
	}

	fmt.Println("🟥 Index created successfully")
}
