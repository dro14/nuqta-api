package yordamchi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dro14/nuqta-service/yordamchi/types"
)

func (y *Yordamchi) send(ctx context.Context, request any) (*http.Response, error) {
	resp, err := y.makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return nil, errors.New(resp.Status)
	default:
		bts, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		response := &types.Response{}
		err = json.Unmarshal(bts, response)
		if err != nil {
			log.Printf("user %s: %s\ncan't decode response: %s\nbody: %s", id(ctx), resp.Status, err, bts)
			return nil, err
		}

		switch {
		case strings.Contains(response.Error.Message, "This model's maximum context length is"):
			err = ErrContextLength
		case strings.Contains(response.Error.Message, "Your request was rejected as a result of our safety system"):
			err = ErrInappropriate
		case strings.Contains(response.Error.Message, "Timeout while downloading"):
			return nil, ErrTimeout
		case strings.Contains(response.Error.Message, "Error while downloading"):
			return nil, ErrDownload
		case resp.StatusCode == http.StatusBadRequest:
			err = ErrBadRequest
		}

		log.Printf("user %s: %s\ntype: %s\nmessage: %s", id(ctx), resp.Status, response.Error.Type, response.Error.Message)
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New(response.Error.Type)
		}
	}
}

func (y *Yordamchi) makeRequest(ctx context.Context, request any) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("user %s: can't encode request: %s", id(ctx), err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, y.endpoint, &buffer)
	if err != nil {
		log.Printf("user %s: can't create request: %s", id(ctx), err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", y.key)

	var client http.Client
	client.Timeout = 10 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
