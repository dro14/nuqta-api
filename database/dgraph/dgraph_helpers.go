package dgraph

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

const retryAttempts = 5

func (d *Dgraph) getJson(ctx context.Context, query string, vars map[string]string) ([]byte, error) {
	var lastErr error
	for i := 0; i < retryAttempts; i++ {
		txn := d.client.NewReadOnlyTxn().BestEffort()
		resp, err := txn.QueryWithVars(ctx, query, vars)
		if err == nil {
			return resp.Json, nil
		}
		lastErr = err
		log.Printf("can't get nquads: %s%s", err, query)
	}
	log.Printf("failed to get nquads after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Dgraph) setJson(ctx context.Context, object any) (*api.Response, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}
	var lastErr error
	for i := 0; i < retryAttempts; i++ {
		resp, err := d.client.NewTxn().Mutate(ctx, mutation)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		log.Printf("can't set json: %s\n%s", err, bytes)
	}
	log.Printf("failed to set json after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Dgraph) deleteJson(ctx context.Context, object any) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	mutation := &api.Mutation{
		DeleteJson: bytes,
		CommitNow:  true,
	}
	var lastErr error
	for i := 0; i < retryAttempts; i++ {
		_, err := d.client.NewTxn().Mutate(ctx, mutation)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("can't delete json: %s\n%s", err, bytes)
	}
	log.Printf("failed to delete json after %d attempts", retryAttempts)
	return lastErr
}
