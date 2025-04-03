package yordamchi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dro14/nuqta-service/yordamchi/types"
)

type key string

var (
	ErrBadRequest    = errors.New("bad request")
	ErrSpit          = errors.New("model spits")
	ErrStream        = errors.New("stream error")
	ErrEmpty         = errors.New("empty response")
	ErrDownload      = errors.New("download error")
	ErrTimeout       = errors.New("download timeout")
	ErrInappropriate = errors.New("inappropriate request")
	ErrContextLength = errors.New("context length exceeded")
)

func sleep(retryDelay *time.Duration) {
	if *retryDelay > 0 {
		log.Print("retrying request after ", *retryDelay)
		time.Sleep(*retryDelay)
		*retryDelay *= 2
	}
}

func id(ctx context.Context) string {
	return ctx.Value(key("firebase_uid")).(string)
}

func getCompletion(response *types.Response) string {
	content, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return ""
	}
	return content
}

func decodeResponse(ctx context.Context, resp *http.Response) (*types.Response, error) {
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("user %s: can't read response: %s", id(ctx), err)
		return nil, err
	}

	response := &types.Response{}
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
