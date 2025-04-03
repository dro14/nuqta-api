package yordamchi

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/nuqta-service/yordamchi/types"
)

func (y *Yordamchi) Completions(ctx context.Context, messages []types.Message) (*types.Response, error) {
	request := &types.Completions{
		Model:       "gpt-4o-mini-2024-07-18",
		Messages:    messages,
		MaxTokens:   3072,
		Temperature: 0.5,
		User:        id(ctx),
	}

	resp, err := y.send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	response, err := decodeResponse(ctx, resp)
	if err != nil {
		return nil, err
	}

	completion := getCompletion(response)
	completion = strings.TrimSpace(completion)
	if len(completion) == 0 {
		return nil, ErrEmpty
	}

	finishReason := response.Choices[0].FinishReason
	if finishReason != "stop" {
		if finishReason != "length" {
			log.Printf("user %s: finish reason is %q", id(ctx), finishReason)
		}
		if didModelSpit(completion) {
			return nil, ErrSpit
		}
	}

	return response, nil
}
