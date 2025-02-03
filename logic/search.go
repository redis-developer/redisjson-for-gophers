package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

const searchQuery = "@actors:{Keanu Reeves} @genres:{action} @rating:[7.0 +inf] @year:[1995 2005]"

func SearchBestMatrixMovies(ctx context.Context, redisClient *redis.Client) {
	searchResult, err := redisClient.FTSearchWithArgs(ctx, IndexName, searchQuery, &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{
			{FieldName: "title", As: "title"},
		},
	}).Result()

	if err != nil {
		log.Printf("Error searching for movies: %v", err)
	}

	if searchResult.Total > 0 {
		fmt.Println("🟥 Document search results:")
		for _, doc := range searchResult.Docs {
			fmt.Printf("   ✅ %s \n", doc.Fields["title"])
		}
	}
}
