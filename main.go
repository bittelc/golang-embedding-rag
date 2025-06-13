package main

import (
	"fmt"
	"golang-embedding-rag/input"
	"golang-embedding-rag/openai"
	"log/slog"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// curl https://api.openai.com/v1/embeddings \
	//  -H "Content-Type: application/json" \
	//  -H "Authorization: Bearer $OPENAI_API_KEY" \
	//  -d '{
	//    "input": "Your text string goes here",
	//    "model": "text-embedding-3-small"
	//  }'
	textToEmbed, err := input.GetUserInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining user input: %v\n", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "OpenAi key not set, required environment variable: $OPENAI_API_KEY")
		os.Exit(1)
	}
	client := openai.NewClient(apiKey)
	req, err := client.CreateRequest(textToEmbed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}
	resp, err := client.CreateEmbeddings(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with response from embeddings request: %v\n", err)
		os.Exit(1)
	}
	spew.Dump(resp)

	slog.Info("Reached end of program")
}
