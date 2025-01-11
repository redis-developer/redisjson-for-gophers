package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"redisjson4gophers/domain"
	"sync"
)

func LoadEmbeddingsFromFile(fileName string) (map[string][]float32, error) {
	const (
		concurrency = 5
	)
	var (
		embeddings = make(map[string][]float32)
		waitGroup  = new(sync.WaitGroup)
		workQueue  = make(chan domain.MovieEmbeddings)
		mutex      = &sync.Mutex{}
		errChan    = make(chan error, concurrency)
	)

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var allEmbeddings []domain.MovieEmbeddings
	err = json.Unmarshal(fileContent, &allEmbeddings)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	for i := 0; i < concurrency; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for embedding := range workQueue {
				mutex.Lock()
				embeddings[embedding.MovieKey] = embedding.PlotEmbeddings
				mutex.Unlock()
			}
		}()
	}

	go func() {
		for _, embedding := range allEmbeddings {
			workQueue <- embedding
		}
		close(workQueue)
	}()

	waitGroup.Wait()
	close(errChan)

	fmt.Printf("🟥 Embeddings loaded from the file: %d \n", len(embeddings))
	return embeddings, nil
}
