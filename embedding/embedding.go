package embedding

import (
	"context"
	"fmt"

	cohere "github.com/cohere-ai/cohere-go/v2"
	"github.com/cohere-ai/cohere-go/v2/client"
)

// Client wraps the Cohere client with custom functionality
type Client struct {
	client *client.Client
}

// NewClient creates a new Cohere client with the provided API key
func NewClient(apiKey string) *Client {
	c := client.NewClient(client.WithToken(apiKey))
	return &Client{
		client: c,
	}
}

// CreateEmbeddings creates embeddings for the given text
func (c *Client) CreateEmbeddings(text string) (*cohere.EmbedByTypeResponse, error) {
	req := cohere.V2EmbedRequest{
		Texts:     []string{text},
		Model:     "embed-english-v3.0",
		InputType: cohere.EmbedInputType("search_query"),
	}

	resp, err := c.client.V2.Embed(context.Background(), &req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	return resp, nil
}
