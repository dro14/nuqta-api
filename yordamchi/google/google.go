package google

import (
	"log"
	"os"
)

type Google struct {
	key      string
	endpoint string
}

func New() *Google {
	key, ok := os.LookupEnv("GOOGLE_API_KEY")
	if !ok {
		log.Fatal("google api key is not specified")
	}

	return &Google{
		key:      key,
		endpoint: "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
	}
}
