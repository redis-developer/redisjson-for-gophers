package logic

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
	"log"
)

func CreateEmbedding(ctx context.Context, text string) []float64 {
	client := openai.NewClient()

	params := openai.EmbeddingNewParams{
		Input:          openai.F[openai.EmbeddingNewParamsInputUnion](shared.UnionString(text)),
		Model:          openai.F(openai.EmbeddingModelTextEmbeddingAda002),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	}

	response, err := client.Embeddings.New(ctx, params)
	if err != nil {
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