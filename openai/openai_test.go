package oai

import (
	"context"
	"errors"
	"fmt"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

// MockOpenAIClient is a mock implementation of the OpenAI client
type MockOpenAIClient struct {
	CreateEmbeddingsFunc func(ctx context.Context, conv openai.EmbeddingRequest) (openai.EmbeddingResponse, error)
}

func (m *MockOpenAIClient) CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequest) (openai.EmbeddingResponse, error) {
	if m.CreateEmbeddingsFunc != nil {
		return m.CreateEmbeddingsFunc(ctx, conv)
	}
	return openai.EmbeddingResponse{}, nil
}

// TestableClient is a version of Client that can be injected with a mock
type TestableClient struct {
	client interface {
		CreateEmbeddings(ctx context.Context, conv openai.EmbeddingRequest) (openai.EmbeddingResponse, error)
	}
}

func (c TestableClient) CreateEmbeddings(req *openai.EmbeddingRequest) (*openai.EmbeddingResponse, error) {
	resp, err := c.client.CreateEmbeddings(context.Background(), *req)
	if err != nil {
		return nil, fmt.Errorf("request CreateEmbeddings failed: %v", err)
	}

	return &resp, nil
}

func TestCreateEmbeddings(t *testing.T) {
	tests := []struct {
		name           string
		request        *openai.EmbeddingRequest
		mockResponse   openai.EmbeddingResponse
		mockError      error
		expectedError  bool
		errorSubstring string
	}{
		{
			name: "successful embedding creation",
			request: &openai.EmbeddingRequest{
				Input: []string{"test text"},
				Model: openai.SmallEmbedding3,
			},
			mockResponse: openai.EmbeddingResponse{
				Object: "list",
				Data: []openai.Embedding{
					{
						Object:    "embedding",
						Index:     0,
						Embedding: []float32{0.1, 0.2, 0.3},
					},
				},
				Model: openai.SmallEmbedding3,
				Usage: openai.Usage{
					PromptTokens: 2,
					TotalTokens:  2,
				},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "API error handling",
			request: &openai.EmbeddingRequest{
				Input: []string{"test text"},
				Model: openai.SmallEmbedding3,
			},
			mockResponse:   openai.EmbeddingResponse{},
			mockError:      errors.New("API quota exceeded"),
			expectedError:  true,
			errorSubstring: "request CreateEmbeddings failed",
		},
		{
			name: "empty input handling",
			request: &openai.EmbeddingRequest{
				Input: []string{},
				Model: openai.SmallEmbedding3,
			},
			mockResponse: openai.EmbeddingResponse{
				Object: "list",
				Data:   []openai.Embedding{},
				Model:  openai.SmallEmbedding3,
				Usage: openai.Usage{
					PromptTokens: 0,
					TotalTokens:  0,
				},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "multiple inputs",
			request: &openai.EmbeddingRequest{
				Input: []string{"first text", "second text"},
				Model: openai.SmallEmbedding3,
			},
			mockResponse: openai.EmbeddingResponse{
				Object: "list",
				Data: []openai.Embedding{
					{
						Object:    "embedding",
						Index:     0,
						Embedding: []float32{0.1, 0.2, 0.3},
					},
					{
						Object:    "embedding",
						Index:     1,
						Embedding: []float32{0.4, 0.5, 0.6},
					},
				},
				Model: openai.SmallEmbedding3,
				Usage: openai.Usage{
					PromptTokens: 4,
					TotalTokens:  4,
				},
			},
			mockError:     nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &MockOpenAIClient{
				CreateEmbeddingsFunc: func(ctx context.Context, conv openai.EmbeddingRequest) (openai.EmbeddingResponse, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			// Create testable client with mock
			client := TestableClient{client: mockClient}

			// Test the function
			result, err := client.CreateEmbeddings(tt.request)

			// Check error expectation
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorSubstring != "" && !contains(err.Error(), tt.errorSubstring) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorSubstring, err.Error())
				}
				if result != nil {
					t.Errorf("Expected nil result on error, got %+v", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Expected non-nil result, got nil")
				}
				if result != nil {
					// Verify response structure
					if result.Object != tt.mockResponse.Object {
						t.Errorf("Expected Object %q, got %q", tt.mockResponse.Object, result.Object)
					}
					if result.Model != tt.mockResponse.Model {
						t.Errorf("Expected Model %q, got %q", tt.mockResponse.Model, result.Model)
					}
					if len(result.Data) != len(tt.mockResponse.Data) {
						t.Errorf("Expected %d embeddings, got %d", len(tt.mockResponse.Data), len(result.Data))
					}
				}
			}
		})
	}
}

func TestCreateEmbeddingsNilRequest(t *testing.T) {
	// Create mock client that should not be called
	mockClient := &MockOpenAIClient{
		CreateEmbeddingsFunc: func(ctx context.Context, conv openai.EmbeddingRequest) (openai.EmbeddingResponse, error) {
			t.Error("CreateEmbeddings should not be called with nil request")
			return openai.EmbeddingResponse{}, nil
		},
	}

	client := TestableClient{client: mockClient}

	// This should panic due to dereferencing nil pointer
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil request")
		}
	}()

	_, _ = client.CreateEmbeddings(nil)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()))
}
