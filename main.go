package main

import (
	"fmt"
	"golang-embedding-rag/openai"
	"log/slog"
	"os"
)

func main() {
	// curl https://api.openai.com/v1/embeddings \
	//  -H "Content-Type: application/json" \
	//  -H "Authorization: Bearer $OPENAI_API_KEY" \
	//  -d '{
	//    "input": "Your text string goes here",
	//    "model": "text-embedding-3-small"
	//  }'
	//
	_, err := openai.CreateRequest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}
	slog.Info("Reached end of program")
}
