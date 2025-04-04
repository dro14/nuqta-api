package yordamchi

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/e"
)

var SystemMessage = map[string]string{
	"oz": `\n\nSen Telegramdagi Yordamchi nomli, matn va rasmlarni tushuna oladigan, xushmuomala chatbotsan. Standart til: O'zbekcha (lotin).`,
	"uz": `\n\nСен Телеграмдаги Ёрдамчи номли, матн ва расмларни тушуна оладиган, хушмуомала чатботсан. Стандарт тил: Ўзбекча (кирил).`,
	"ru": `\n\nТы являешься дружелюбным чатботом в Телеграме под названием Yordamchi, который понимает текст и изображения. Язык по умолчанию: Русский.`,
	"en": `\n\nYour are a friendly chatbot in Telegram called Yordamchi, that can follow text and images. Default language: English.`,
}

func (y *Yordamchi) Respond(ctx context.Context, conversation []string, language, provider string) (string, error) {
	if len(conversation) > 3 {
		conversation = conversation[len(conversation)-3:]
	}
	now := time.Now().Format(time.DateTime)
	conversation = append(
		[]string{now + SystemMessage[language]},
		conversation...,
	)

	retryDelay := 5 * time.Second
	attempts := 0
	errMsg := ""
Retry:
	attempts++
	var response string
	var err error
	switch provider {
	case "openai":
		response, err = y.openai.Completions(ctx, conversation)
	case "google":
		response, err = y.google.GenerateContent(ctx, conversation)
	}
	if err != nil {
		errMsg = err.Error()

		switch {
		case errors.Is(err, e.ErrContextLength),
			errors.Is(err, e.ErrInappropriate),
			errors.Is(err, e.ErrBadRequest):
			return "", err
		case errors.Is(err, e.ErrStream), errors.Is(err, e.ErrSpit):
			if attempts > 1 {
				log.Printf("user %s: %q failed after %d attempts", id(ctx), errMsg, attempts)
				return "", err
			}
			retryDelay = 0
		case errors.Is(err, e.ErrDownload):
			if len(conversation) > 2 {
				conversation = append(conversation[:1], conversation[3:]...)
				attempts--
			}
			retryDelay = 0
		case errors.Is(err, e.ErrTimeout), strings.Contains(errMsg, "context deadline exceeded"),
			errMsg == "500 Internal Server Error", errMsg == "502 Bad Gateway":
			retryDelay = 0
		}

		if attempts < 10 {
			sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("user %s: %q failed after %d attempts", id(ctx), errMsg, attempts)
			return "", err
		}
	}

	if attempts > 1 {
		log.Printf("user %s: %q was handled after %d attempts", id(ctx), errMsg, attempts)
	}

	return response, nil
}
