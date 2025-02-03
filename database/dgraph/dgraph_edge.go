package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/e"
)

var edges = []string{
	"follow",
	"like",
	"repost",
	"click",
	"viewed_by",
}

func (d *Dgraph) CreateEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> <%s> (timestamp=%d) .", source, edge, target, time.Now().Unix()))
	mutation := &api.Mutation{SetNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	return err
}

func (d *Dgraph) DeleteEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> <%s> .", source, edge, target))
	mutation := &api.Mutation{DelNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	return err
}

func (d *Dgraph) doesEdgeExist(ctx context.Context, source, edge, target string) (bool, error) {
	if !slices.Contains(edges, edge) {
		return false, e.ErrUnknownEdge
	}
	query := fmt.Sprintf(edgeQuery, source, edge, target)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return false, err
	}

	var response map[string][]any
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return false, err
	}

	return len(response["edges"]) > 0, nil
}
