package openai

import (
	"fmt"
	"net/http"
	"net/url"
)

const openAiUrl = "https://api.openai.com/v1/embeddings"

func CreateRequest() (http.Request, error) {
	var openAiReq http.Request
	url, err := url.Parse(openAiUrl)
	if err != nil {
		return openAiReq, fmt.Errorf("couldn't create Url for OpenAI endpoint", err)
	}
	req := http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			"Authorization": []string{"Bearer YOUR_API_KEY"},
		},
	}
	return req, nil
}
