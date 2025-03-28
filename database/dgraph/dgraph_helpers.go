package dgraph

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

const retryAttempts = 5

func (d *Dgraph) get(ctx context.Context, query string, vars map[string]string) ([]byte, error) {
	var lastErr error
	for range retryAttempts {
		txn := d.client.NewReadOnlyTxn().BestEffort()
		resp, err := txn.QueryWithVars(ctx, query, vars)
		if err == nil {
			return resp.Json, nil
		}
		lastErr = err
		log.Printf("can't get: %s%s", err, query)
	}
	log.Printf("failed to get after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Dgraph) set(ctx context.Context, object any) (*api.Response, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	log.Printf("%s\n", bytes)
	mutation := &api.Mutation{
		SetJson:   bytes,
		CommitNow: true,
	}
	var lastErr error
	for range retryAttempts {
		txn := d.client.NewTxn()
		resp, err := txn.Mutate(ctx, mutation)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		log.Printf("can't set: %s\n%s", err, bytes)
	}
	log.Printf("failed to set after %d attempts", retryAttempts)
	return nil, lastErr
}

func (d *Dgraph) delete(ctx context.Context, object any) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	mutation := &api.Mutation{
		DeleteJson: bytes,
		CommitNow:  true,
	}
	var lastErr error
	for range retryAttempts {
		txn := d.client.NewTxn()
		_, err := txn.Mutate(ctx, mutation)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("can't delete: %s\n%s", err, bytes)
	}
	log.Printf("failed to delete after %d attempts", retryAttempts)
	return lastErr
}
