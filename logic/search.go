package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const searchQuery = "(*)=>[KNN 1 @plotEmbeddings $vector]"

func SearchMovieWithVectorField(ctx context.Context, redisClient *redis.Client) {
	//query := "He seeks revenge for the death of his family. He is a vigilante."
	query := "He is also known as the man without fear."

	rawResult, err := redisClient.FTSearchWithArgs(ctx, IndexName, searchQuery, &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{
			{FieldName: "$.title", As: "title"},
			{FieldName: "$.plot", As: "plot"},
		},
		Params:         map[string]interface{}{"vector": ConvertFloatsToByte(CreateEmbedding(ctx, query))},
		DialectVersion: 2,
	}).RawResult()
	if err != nil {
		log.Printf("Error executing the search: %v", err)
		return
	}

	if rawResult != nil {
		rawResults := rawResult.(map[interface{}]interface{})

		var movieTitle string
		var moviePlot string
		if rawResults["total_results"].(int64) > 0 {
			results := rawResults["results"].([]interface{})
			for _, result := range results {
				movie := result.(map[interface{}]interface{})["extra_attributes"].(map[interface{}]interface{})
				movieTitle = movie["title"].(string)
				moviePlot = movie["plot"].(string)
			}
		}
		fmt.Println("🟥 Similarity search result: ")
		fmt.Printf("   %s \n   %s \n", movieTitle, moviePlot)
	}
}
