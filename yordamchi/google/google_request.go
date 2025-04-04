package google

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (g *Google) send(ctx context.Context, request any) (*http.Response, error) {
	resp, err := g.makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("user %s: can't read response: %s", id(ctx), err)
		return nil, err
	}
	_ = resp.Body.Close()
	log.Printf("user %s: %s\nbody: %s", id(ctx), resp.Status, bts)
	return nil, errors.New(resp.Status)
}

func (g *Google) makeRequest(ctx context.Context, request any) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("user %s: can't encode request: %s", id(ctx), err)
		return nil, err
	}

	url := fmt.Sprintf(g.endpoint, ctx.Value("model"), g.key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buffer)
	if err != nil {
		log.Printf("user %s: can't create request: %s", id(ctx), err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var client http.Client
	client.Timeout = 10 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
