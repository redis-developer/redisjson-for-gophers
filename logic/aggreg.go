package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func MovieCountPerGenreAgg(ctx context.Context, redisClient *redis.Client) {
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

	aggregResult, err := redisClient.FTAggregateWithArgs(ctx, IndexName, "*", aggregOptions).Result()
	if err != nil {
		log.Printf("Error executing the aggregation: %v", err)
		return
	}

	if aggregResult.Total > 0 {
		fmt.Printf("🟥 Top 5 Genres and their Movie Count: \n")
		for i, aggregRow := range aggregResult.Rows {
			fmt.Printf("   %d)️ %s = %s\n", i+1, aggregRow.Fields["genres"], aggregRow.Fields["Count"])
		}
	}
}
