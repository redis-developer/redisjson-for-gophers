package logic

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/openai/openai-go"
	"log"
)

func CreateEmbedding(ctx context.Context, text string) []float64 {
	client := openai.NewClient()

	response, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: []string{text},
		},
		Model:          openai.EmbeddingModelTextEmbeddingAda002,
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	})
	if err != nil {
		log.Printf("Error creating embedding: %v", err)
		return nil
	}

	return response.Data[0].Embedding
}

func ConvertFloatsToByte(floats []float64) []byte {
	if len(floats) != 1536 {
		log.Fatalf("Expected 1536 dimensions, but got %d", len(floats))
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, floats)
	if err != nil {
		log.Fatalf("binary.Write failed: %v", err)
	}

	result := buf.Bytes()
	return result
}
