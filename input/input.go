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
