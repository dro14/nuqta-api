package openai

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/nuqta-service/e"
)

func (o *OpenAI) Completions(ctx context.Context, conversation []string) (string, error) {
	var messages []Message
	for i, text := range conversation {
		var role string
		if i == 0 {
			role = "system"
		} else if i%2 != 0 {
			role = "user"
		} else {
			role = "assistant"
		}
		messages = append(messages, Message{
			Content: text,
			Role:    role,
		})
	}

	request := &Completions{
		Model:       "gpt-4o-mini-2024-07-18",
		Messages:    messages,
		MaxTokens:   3072,
		Temperature: 0.5,
		User:        id(ctx),
	}

	resp, err := o.send(ctx, request)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	response, err := decodeResponse(ctx, resp)
	if err != nil {
		return "", err
	}

	completion := getCompletion(response)
	completion = strings.TrimSpace(completion)
	if len(completion) == 0 {
		return "", e.ErrEmpty
	}

	finishReason := response.Choices[0].FinishReason
	if finishReason != "stop" {
		if finishReason != "length" {
			log.Printf("user %s: finish reason is %q", id(ctx), finishReason)
		}
		if didModelSpit(completion) {
			return "", e.ErrSpit
		}
	}

	return completion, nil
}
