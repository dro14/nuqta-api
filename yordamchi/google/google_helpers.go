package google

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
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return ""
	}
	return response.Candidates[0].Content.Parts[0].Text
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
