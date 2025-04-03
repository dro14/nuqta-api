package yordamchi

import (
	"log"
	"os"
)

type Yordamchi struct {
	key      string
	endpoint string
}

func New() *Yordamchi {
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		log.Fatal("openai api key is not specified")
	}

	return &Yordamchi{
		key:      key,
		endpoint: "https://api.openai.com/v1/chat/completions",
	}
}
