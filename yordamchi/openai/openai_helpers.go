package openai

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func id(ctx context.Context) string {
	return ctx.Value("firebase_uid").(string)
}

func getCompletion(response *Response) string {
	content, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return ""
	}
	return content
}

func decodeResponse(ctx context.Context, resp *http.Response) (*Response, error) {
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("user %s: can't read response: %s", id(ctx), err)
		return nil, err
	}

	response := &Response{}
	err = json.Unmarshal(bts, response)
	if err != nil {
		log.Printf("user %s: can't decode response: %s\nbody: %s", id(ctx), err, bts)
		return nil, err
	}
	return response, nil
}

func didModelSpit(completion string) bool {
	runes := []rune(completion)
	Length := len(runes)
	maxSubLength := Length / 10
	for SubLength := 2; SubLength <= maxSubLength; SubLength++ {
		i := 0
		for i <= Length-SubLength {
			substring := string(runes[i : i+SubLength])
			count := 1
			for i+SubLength < Length && string(runes[i+SubLength:min(i+2*SubLength, cap(runes))]) == substring {
				count++
				i += SubLength
				if count > 10 {
					return true
				}
			}
			i += SubLength
		}
	}
	return false
}
