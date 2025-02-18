package dgraph

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v240/protos/api"
)

const retryAttempts = 5

func (d *Dgraph) get(ctx context.Context, query string, vars map[string]string) ([]byte, error) {
	var err error
	for i := 0; i < retryAttempts; i++ {
		resp, err := d.client.NewTxn().QueryWithVars(ctx, query, vars)
		if err == nil {
			return resp.Json, nil
		}
		log.Printf("can't get nquads: %s\n%s", err, query)
	}
	log.Printf("failed to get nquads after %d attempts", retryAttempts)
	return nil, err
}

func (d *Dgraph) setNquads(ctx context.Context, query string) error {
	mutation := &api.Mutation{
		SetNquads: []byte(query),
		CommitNow: true,
	}
	var err error
	for i := 0; i < retryAttempts; i++ {
		_, err = d.client.NewTxn().Mutate(ctx, mutation)
		if err == nil {
			return nil
		}
		log.Printf("can't set nquads: %s\n%s", err, query)
	}
	log.Printf("failed to set nquads after %d attempts", retryAttempts)
	return err
}

func (d *Dgraph) deleteNquads(ctx context.Context, query string) error {
	mutation := &api.Mutation{
		DelNquads: []byte(query),
		CommitNow: true,
	}
	var err error
	for i := 0; i < retryAttempts; i++ {
		_, err = d.client.NewTxn().Mutate(ctx, mutation)
		if err == nil {
			return nil
		}
		log.Printf("can't delete nquads: %s\n%s", err, query)
	}
	log.Printf("failed to delete nquads after %d attempts", retryAttempts)
	return err
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
	var resp *api.Response
	for i := 0; i < retryAttempts; i++ {
		resp, err = d.client.NewTxn().Mutate(ctx, mutation)
		if err == nil {
			return resp, nil
		}
		log.Printf("can't set json: %s\n%s", err, bytes)
	}
	log.Printf("failed to set json after %d attempts", retryAttempts)
	return nil, err
}
