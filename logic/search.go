package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const searchQuery = "(*)=>[KNN 1 @plotEmbedding $vector]"

func SearchMovieWithVectorField(ctx context.Context, redisClient *redis.Client) {
	queryParam := "He wears a skull in his chest and seeks revenge for his family."

	searchResult, err := redisClient.FTSearchWithArgs(ctx, IndexName, searchQuery, &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{
			{FieldName: "$.title", As: "title"},
			{FieldName: "$.plot", As: "plot"},
		},
		Params:         map[string]interface{}{"vector": ConvertFloatsToByte(CreateEmbedding(ctx, queryParam))},
		DialectVersion: 2,
	}).Result()

	if err != nil {
		log.Printf("Error searching for movie: %v", err)
	}

	if searchResult.Total > 0 {
		fmt.Println("🟥 Similarity search result:")
		doc := searchResult.Docs[0]
		fmt.Printf("   🎥 %s \n", doc.Fields["title"])
		fmt.Printf("   💬 %s \n", doc.Fields["plot"])
	}
}
