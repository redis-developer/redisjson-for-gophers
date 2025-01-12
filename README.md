# Redis JSON for Gophers

This project contains an example that showcases different features from the official [Go Client for Redis](https://github.com/redis/go-redis) that you can use as a reference about how to get started with the support for JSON in Redis in your Go apps. It is not intended to provide the full spectrum of what the client is capable of—but it certainly puts you on the right track.

You can run this code with an Redis instance running locally, to which you can leverage the [Docker Compose code](./docker-compose.yml) available in the project. Alternatively, you can also run this code with [Redis Cloud](https://redis.io/cloud/) that can be easily created using the [Terraform code](./redis-cloud.tf) also available in the project.

## Examples available in this project:

### 🟥 Movies Loading

The data model from this project is a collection of movies from the file [movies.json](./movies.json). This file will be [loaded](logic/movies.go) in memory and made available within the context, which the other functions will work with. Here is an example of a movie:

```json
{
  "title": "Blade",
  "year": 1998,
  "plot": "A half-vampire, half-mortal man becomes a protector of the mortal race, while slaying evil vampires.",
  "runningTime": 7200,
  "releaseDate": "1998-08-19T00:00:00Z",
  "rating": 7,
  "genres": [
    "Action",
    "Fantasy",
    "Horror"
  ],
  "actors": [
    "Wesley Snipes",
    "Stephen Dorff",
    "Kris Kristofferson"
  ],
  "directors": [
    "Stephen Norrington"
  ]
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

### 🟥 Aggregation Analytics

The project also runs a very interesting [aggregation](logic/aggreg.go) to find out the top five genres and their respective movie counts. Just like you would do with:

```bash
FT.AGGREGATE movies_index * GROUPBY 1 @genres REDUCE COUNT 0 AS Count SORTBY 2 @Count DESC MAX 5
```

### 🟥 Semantic Searches

Obviously, this project couldn't leave behind an example of a vector search. The implemented [search](logic/search.go) to find a movie given a plot similar to the one you are looking for. 

```bash
FT.SEARCH movies_index "(*)=>[KNN 1 @plotEmbedding $vector]" RETURN 6 $.title AS title $.plot AS plot PARAMS 2 vector "<QUERY_PARAM_EMBEDDING>" DIALECT 2
```

# License

This project is licensed under the [Apache 2.0 License](./LICENSE).
