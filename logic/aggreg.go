package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"redisjson4gophers/domain"
)

func MovieCountPerGenreAgg(ctx context.Context) {
	redisClient := ctx.Value(domain.ClientKey).(*redis.Client)

	aggregOptions := &redis.FTAggregateOptions{
		GroupBy: []redis.FTAggregateGroupBy{
			{
				Fields: []interface{}{"@genres"},
				Reduce: []redis.FTAggregateReducer{
					{
						Reducer: redis.SearchCount,
						As:      "Count",
					},
				},
			},
		},
		SortBy: []redis.FTAggregateSortBy{
			{
				FieldName: "@Count",
				Desc:      true,
			},
		},
		SortByMax: 5,
	}

	rawResult, err := redisClient.FTAggregateWithArgs(ctx, indexName, "*", aggregOptions).RawResult()
	if err != nil {
		log.Printf("Error executing the aggregation: %v", err)
		return
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
