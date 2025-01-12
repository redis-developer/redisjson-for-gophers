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

	rawResult, err := redisClient.FTSearchWithArgs(ctx, IndexName, searchQuery, &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{
			{FieldName: "$.title", As: "title"},
			{FieldName: "$.plot", As: "plot"},
		},
		Params:         map[string]interface{}{"vector": ConvertFloatsToByte(CreateEmbedding(ctx, queryParam))},
		DialectVersion: 2,
	}).RawResult()
	if err != nil {
		log.Printf("Error executing the search: %v", err)
		return
	}

	if rawResult != nil {
		rawResults := rawResult.(map[interface{}]interface{})

		if rawResults["total_results"].(int64) > 0 {
			results := rawResults["results"].([]interface{})
			movie := results[0].(map[interface{}]interface{})["extra_attributes"].(map[interface{}]interface{})
			fmt.Println("🟥 Similarity search result:")
			fmt.Printf("   🎥 %s \n", movie["title"].(string))
			fmt.Printf("   💬 %s \n", movie["plot"].(string))
		}
	}
}
