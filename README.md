# Redis JSON for Gophers

This project contains an example that showcases different features from the official [Go Client for Redis](https://github.com/redis/go-redis) that you can use as a reference about how to get started with the support for JSON in Redis in your Go apps. It is not intended to provide the full spectrum of what the client is capable of—but it certainly puts you on the right track.

You can run this code with an Redis instance running locally, to which you can leverage the [Docker Compose code](./docker-compose.yml) available in the project. Alternatively, you can also run this code with [Redis Cloud](https://redis.io/cloud/) that can be easily created using the [Terraform code](./redis-cloud.tf) also available in the project.

## Examples available in this project:

### 🟥 Movies Loading

The data model from this project is a collection of movies from the file [movies.json](./movies.json). This file will be [loaded](logic/movies.go) in memory and made available within the context, which the other functions will work with. Here is an example of a movie:

```json
{
    "year": 2012,
    "title": "The Avengers",
    "info": {
        "directors": [
            "Joss Whedon"
        ],
        "release_date": "2012-04-11T00:00:00Z",
        "rating": 8.2,
        "genres": [
            "Action",
            "Fantasy"
        ],
        "image_url": "http://ia.media-imdb.com/images/M/MV5BMTk2NTI1MTU4N15BMl5BanBnXkFtZTcwODg0OTY0Nw@@._V1_SX400_.jpg",
        "plot": "Nick Fury of S.H.I.E.L.D. assembles a team of superhumans to save the planet from Loki and his army.",
        "rank": 48,
        "running_time_secs": 8580,
        "actors": [
            "Robert Downey Jr.",
            "Chris Evans",
            "Scarlett Johansson"
        ]
    }
}
```

### 🟥 Connection Handling

Once the movies are loaded, the code will create a [connection](logic/connect.go) with Redis and make this connection available within the context as well.

```go
connOpts, err := redis.ParseURL(redisConnectionURL)
if err != nil {
    panic(fmt.Errorf("failed to parse Redis URL: %w", err))
}
redisClient := redis.NewClient(connOpts)

_, err = redisClient.Ping(ctx).Result()
if err != nil {
panic(fmt.Errorf("error connecting with Redis: %w", err))
}

return context.WithValue(ctx, domain.ClientKey, redisClient)
```

### 🟥 Document Indexing

All the movies will be [indexed](logic/index.go) in Redis. The example uses [Redis Pipelining](https://redis.io/docs/latest/develop/use/pipelining/) to index documents.

```go
pipeline := redisClient.Pipeline()
for movieID, movie := range movies {
    movieAsJSON, err := json.Marshal(movie)
    if err != nil {
        log.Printf("Error marshaling movie into JSON: %v", err)
    }
    pipeline.JSONSet(ctx, keyPrefix+strconv.Itoa(movieID), "$", string(movieAsJSON))
}

_, _ := pipeline.Exec(ctx)
```

### 🟥 Document Lookup

An example of [document lookup](logic/lookup.go) is also available. Out of all movies loaded, an key will be randomly selected, and the document associated with this key will be looked up. Just like you would do with:

```bash
JSON.GET movie:1234 $.title
```

### ✅ Doing Searches

Obviously, this project couldn't leave behind an example of a search. The implemented [search](logic/search.go) look for all the best action movies from [Keanu Reeves](https://en.wikipedia.org/wiki/Keanu_Reeves) from 1995 to 2005. This search is the equivalent to:

```bash
FT.SEARCH json_movies_index "@actors:{Keanu Reeves} @genres:{action} @rating:[7.0 +inf] @year:[1995 2005]" RETURN 3 $.title $.year $.rating
```

### ✅ Aggregation Analytics

Finally, the project also runs a very interesting [aggregation](logic/aggreg.go) to find out the top five genres and their respective movie counts. Just like you would do with:

```bash
FT.AGGREGATE json_movies_index * GROUPBY 1 @genres REDUCE COUNT 0 AS Count SORTBY 2 @Count DESC MAX 5
```

# License

This project is licensed under the [Apache 2.0 License](./LICENSE).
