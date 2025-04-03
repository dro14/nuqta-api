package yordamchi

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/yordamchi/types"
)

var SystemMessage = map[string]string{
	"oz": `Sen Telegramdagi Yordamchi nomli, matn va rasmlarni tushuna oladigan, xushmuomala chatbotsan. Standart til: O'zbekcha (lotin). Hozirgi vaqt: `,
	"uz": `Сен Телеграмдаги Ёрдамчи номли, матн ва расмларни тушуна оладиган, хушмуомала чатботсан. Стандарт тил: Ўзбекча (кирил). Ҳозирги вақт: `,
	"ru": `Ты являешься дружелюбным чатботом в Телеграме под названием Yordamchi, который понимает текст и изображения. Язык по умолчанию: Русский. Текущее время: `,
	"en": `Your are a friendly chatbot in Telegram called Yordamchi, that can follow text and images. Default language: English. Current time: `,
}

func (y *Yordamchi) ProcessCompletions(ctx context.Context, conversation []string, language string, firebaseUid string) (string, error) {
	ctx = context.WithValue(ctx, key("firebase_uid"), firebaseUid)
	messages := []types.Message{{
		Role:    "system",
		Content: SystemMessage[language] + time.Now().Format(time.DateTime),
	}}
	if len(conversation) > 3 {
		conversation = conversation[len(conversation)-3:]
	}
	for i, message := range conversation {
		role := "user"
		if i == 1 {
			role = "assistant"
		}
		messages = append(messages, types.Message{
			Role:    role,
			Content: message,
		})
	}

	retryDelay := 5 * time.Second
	attempts := 0
	errMsg := ""
Retry:
	attempts++
	response, err := y.Completions(ctx, messages)
	if err != nil {
		errMsg = err.Error()

		switch {
		case errors.Is(err, ErrContextLength),
			errors.Is(err, ErrInappropriate),
			errors.Is(err, ErrBadRequest):
			return "", err
		case errors.Is(err, ErrStream), errors.Is(err, ErrSpit):
			if attempts > 1 {
				log.Printf("user %s: %q failed after %d attempts", id(ctx), errMsg, attempts)
				return "", err
			}
			retryDelay = 0
		case errors.Is(err, ErrDownload):
			if len(messages) > 2 {
				messages = append(messages[:1], messages[3:]...)
				attempts--
			}
			retryDelay = 0
		case errors.Is(err, ErrTimeout), strings.Contains(errMsg, "context deadline exceeded"),
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

	return getCompletion(response), nil
}
