package logic

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"redisjson4gophers/domain"
)

func MovieCountPerGenreAgg(ctx context.Context) {

	redisClient := ctx.Value(domain.ClientKey).(*redis.Client)

	cmdResult := redisClient.Do(ctx,
		"FT.AGGREGATE", indexName, "*", "GROUPBY", "1",
		"@genres", "REDUCE", "COUNT", "0", "AS", "Count", "SORTBY",
		"2", "@Count", "DESC", "MAX", "5",
	)
	rawResult, err := cmdResult.Result()
	if err != nil {
		log.Printf("Error executing the aggregation: %v", err)
	}

	if rawResult != nil {
		aggregationResults := rawResult.(map[interface{}]interface{})["results"].([]interface{})
		if len(aggregationResults) > 0 {
			fmt.Printf("🟥 Top 5 Genres and their Movie Count: \n")
			for i, aggregResult := range aggregationResults {
				entry := aggregResult.(map[interface{}]interface{})
				extraAttribs := entry["extra_attributes"].(map[interface{}]interface{})
				fmt.Printf("   %d)️ %s = %s\n", i+1, extraAttribs["genres"], extraAttribs["Count"])
			}
		}
	}
}
