package main

import (
	"fmt"
	"golang-embedding-rag/input"
	oai "golang-embedding-rag/openai"
	"log/slog"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "OpenAi key not set, required environment variable: $OPENAI_API_KEY")
		os.Exit(1)
	}
	client := oai.NewClient(apiKey)

	textToEmbed, err := input.GetFileInput("data/input_text") // Can be switched out for GetUserInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting text from file: %v\n", err)
		os.Exit(1)
	}

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
	timestamp := time.Now().Format("20060102.1504")
	filename := fmt.Sprintf("data/output_%s", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	spew.Fdump(file, resp)

	slog.Info("Reached end of program")
}
