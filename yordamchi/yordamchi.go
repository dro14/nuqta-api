package yordamchi

import (
	"github.com/dro14/nuqta-service/yordamchi/google"
	"github.com/dro14/nuqta-service/yordamchi/openai"
)

type Yordamchi struct {
	google *google.Google
	openai *openai.OpenAI
}

func New() *Yordamchi {
	return &Yordamchi{
		google: google.New(),
		openai: openai.New(),
	}
}
