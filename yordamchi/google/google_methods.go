package google

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/nuqta-service/e"
)

func (g *Google) GenerateContent(ctx context.Context, conversation []string) (string, error) {
	request := &Request{
		SystemInstruction: &Content{
			Parts: []Part{{Text: conversation[0]}},
		},
		GenerationConfig: &GenerationConfig{
			MaxOutputTokens: 3072,
			Temperature:     0.5,
		},
	}
	for i, text := range conversation[1:] {
		var role string
		if i%2 == 0 {
			role = "user"
		} else {
			role = "model"
		}
		request.Contents = append(request.Contents, Content{
			Parts: []Part{{Text: text}},
			Role:  role,
		})
	}

	ctx = context.WithValue(ctx, "model", "gemini-2.0-flash-001")
	resp, err := g.send(ctx, request)
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

	finishReason := response.Candidates[0].FinishReason
	if finishReason != "STOP" {
		log.Printf("user %s: finish reason is %q", id(ctx), finishReason)
	}

	return completion, nil
}
