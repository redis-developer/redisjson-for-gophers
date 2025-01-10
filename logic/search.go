package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"redisjson4gophers/domain"
	"strings"
)

const searchQuery = "@actors:{Keanu Reeves} @genres:{action} @rating:[7.0 +inf] @year:[1995 2005]"

func SearchBestMatrixMovies(ctx context.Context) {

	redisClient := ctx.Value(domain.ClientKey).(*redis.Client)

	searchResult := redisClient.FTSearchWithArgs(ctx, indexName, searchQuery, &redis.FTSearchOptions{
		Return: []redis.FTSearchReturn{
			{FieldName: "title", As: "title"},
			{FieldName: "year", As: "year"},
			{FieldName: "rating", As: "rating"},
		},
	})

	if searchResult.RawVal() != nil {
		rawResults := searchResult.RawVal().(map[interface{}]interface{})

		var movieTitles []string
		if rawResults["total_results"].(int64) > 0 {
			results := rawResults["results"].([]interface{})
			for _, result := range results {
				movie := result.(map[interface{}]interface{})["extra_attributes"].(map[interface{}]interface{})
				movieTitles = append(movieTitles, movie["title"].(string))
			}
		}
		fmt.Printf("🟥 Best Matrix movies with Keanu Reeves: [%s] \n", strings.Join(movieTitles, ", "))
	}
}
