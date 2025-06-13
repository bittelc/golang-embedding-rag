package openai

import (
	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	openai.Client
}

func NewClient(key string) Client {
	oClient := openai.NewClient(key)
	return Client{*oClient}
}

func (c *Client) CreateRequest(text string) (*openai.EmbeddingRequest, error) {
	openAiReq := &openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	}
	return openAiReq, nil
}

func (c *Client) CreateEmbeddings(req *openai.EmbeddingRequest) (*openai.EmbeddingResponse, error) {
	return nil, nil
}
