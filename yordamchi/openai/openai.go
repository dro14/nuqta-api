package openai

import (
	"log"
	"os"
)

type OpenAI struct {
	key      string
	endpoint string
}

func New() *OpenAI {
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		log.Fatal("openai api key is not specified")
	}

	return &OpenAI{
		key:      "Bearer " + key,
		endpoint: "https://api.openai.com/v1/chat/completions",
	}
}
