package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type embeddingRequest struct {
	InputText string `json:"inputText"`
}

type embeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GetEmbedding(input string) ([]float32, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	reqBody := embeddingRequest{
		InputText: input,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	resp, err := client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		Body:        jsonBody,
		ModelId:     aws.String("amazon.titan-embed-text-v1"),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return nil, fmt.Errorf("error invoking model: %v", err)
	}

	var embeddingResp embeddingResponse
	err = json.Unmarshal(resp.Body, &embeddingResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	return embeddingResp.Embedding, nil
}
