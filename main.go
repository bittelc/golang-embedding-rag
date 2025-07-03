package main

import (
	"fmt"
	"golang-embedding-rag/embedding"
	"golang-embedding-rag/input"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	apiKey := os.Getenv("COHERE_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Cohere API key not set, required environment variable: $COHERE_API_KEY")
		os.Exit(1)
	}
	client := embedding.NewClient(apiKey)

	// textToEmbed, err := input.GetFileInput("data/input_text") // Can be switched out for GetUserInput()
	textToEmbed, err := input.GetUserInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error extracting text from file: %v\n", err)
		os.Exit(1)
	}

	dataset, err := client.CreateEmbeddings(textToEmbed)
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
	fmt.Println(dataset.GetEmbeddings())

	spew.Fdump(file, dataset)
}
