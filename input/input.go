package input

import (
	"bufio"
	"fmt"
	"os"
)

func GetUserInput() (string, error) {
	var userInput string
	fmt.Print("Text to embed: ")
	reader := bufio.NewReader(os.Stdin)

	userInput, err := reader.ReadString('\n')
	if err != nil {
		return userInput, fmt.Errorf("couldn't create Url for OpenAI endpoint: %v", err)
	}
	return userInput, nil
}

// GetFileInput reads the entire contents of a file at the specified path and returns it as a string.
// Reads all lines concatenating them together, and returns the combined text.
func GetFileInput(path string) (string, error) {
	var prompt string
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("couldn't open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		prompt += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return prompt, err
}
