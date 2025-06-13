package oai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	*openai.Client
}

// NewClient creates a new OpenAI client.
func NewClient(key string) Client {
	oClient := openai.NewClient(key)
	return Client{oClient}
}

// CreateRequest creates a new embedding request for a short string of text
// Should be used for local testing of embedding text. Not intended for use with document-length
// texts.
func (c Client) CreateRequest(text string) (*openai.EmbeddingRequest, error) {
	openAiReq := &openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3,
	}
	return openAiReq, nil
}

// CreateEmbeddings creates a new embedding request for a given request
// Returns an embedding response, which contains the vectorized representation of the request's
// input text
func (c Client) CreateEmbeddings(req *openai.EmbeddingRequest) (*openai.EmbeddingResponse, error) {
	resp, err := c.Client.CreateEmbeddings(context.Background(), *req)
	if err != nil {
		return nil, fmt.Errorf("request CreateEmbeddings failed: %v", err)
	}

	return &resp, nil
}
