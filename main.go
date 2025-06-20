package main

import (
	"fmt"
	"golang-embedding-rag/input"
	oai "golang-embedding-rag/openai"
	"log/slog"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "OpenAi key not set, required environment variable: $OPENAI_API_KEY")
		os.Exit(1)
	}
	client := oai.NewClient(apiKey)

	// textToEmbed, err := input.GetUserInput()
	textToEmbed, err := input.GetFileInput("data/input_text")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining user input: %v\n", err)
		os.Exit(1)
	}

	req, err := client.CreateRequest(textToEmbed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}
	spew.Dump(req)
	resp, err := client.CreateEmbeddings(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error with response from embeddings request: %v\n", err)
		os.Exit(1)
	}
	spew.Dump(resp)

	slog.Info("Reached end of program")
}
