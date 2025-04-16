package yordamchi

import (
	"context"
	"errors"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (y *Yordamchi) Respond(ctx context.Context, provider string, conversation []string) (*models.Message, error) {
	retryDelay := 5 * time.Second
	attempts := 0
	errMsg := ""
Retry:
	attempts++
	var response string
	var err error
	switch provider {
	case "google":
		ctx = context.WithValue(ctx, "model", "gemini-2.0-flash-001")
		response, err = y.google.GenerateContent(ctx, conversation)
	case "openai":
		ctx = context.WithValue(ctx, "model", "gpt-4o-mini-2024-07-18")
		response, err = y.openai.Completions(ctx, conversation)
	default:
		return nil, e.ErrInvalidParams
	}
	if err != nil {
		errMsg = err.Error()

		switch {
		case errors.Is(err, e.ErrContextLength),
			errors.Is(err, e.ErrInappropriate),
			errors.Is(err, e.ErrBadRequest):
			return nil, err
		case errors.Is(err, e.ErrStream), errors.Is(err, e.ErrSpit):
			if attempts > 1 {
				log.Printf("user %s: %q failed after %d attempts", id(ctx), errMsg, attempts)
				return nil, err
			}
			retryDelay = 0
		case errors.Is(err, e.ErrDownload):
			if len(conversation) > 2 {
				conversation = slices.Delete(conversation, 1, 3)
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
			return nil, err
		}
	}

	if attempts > 1 {
		log.Printf("user %s: %q was handled after %d attempts", id(ctx), errMsg, attempts)
	}

	return &models.Message{
		AuthorUid: ctx.Value("model").(string),
		Text:      response,
	}, nil
}
