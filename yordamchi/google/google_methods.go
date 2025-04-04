package google

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/nuqta-service/e"
)

func (g *Google) GenerateContent(ctx context.Context, conversation []string) (string, error) {
	request := &Request{}
	for i, text := range conversation {
		if i == 0 {
			request.SystemInstruction = Content{
				Parts: []Part{{Text: text}},
			}
		} else if i%2 != 0 {
			request.Contents = append(request.Contents, Content{
				Parts: []Part{{Text: text}},
				Role:  "user",
			})
		} else {
			request.Contents = append(request.Contents, Content{
				Parts: []Part{{Text: text}},
				Role:  "model",
			})
		}
	}

	log.Printf("request:\n%+v", *request)

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
		if didModelSpit(completion) {
			return "", e.ErrSpit
		}
	}

	return completion, nil
}
