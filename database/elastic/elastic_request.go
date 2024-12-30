package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (e *Elastic) makeRequest(ctx context.Context, request any, endpoint, method string) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("user %s: can't encode request: %s", id(ctx), err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, e.baseURI+endpoint, &buffer)
	if err != nil {
		log.Printf("user %s: can't create request: %s", id(ctx), err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(e.username, e.password)

	client := &http.Client{}
	client.Timeout = 1 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
