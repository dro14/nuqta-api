package dgraph

import (
	"context"
	"fmt"
	"slices"

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

func (d *Dgraph) CreateEdge(ctx context.Context, node1, edge, node2 string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> <%s> .", node1, edge, node2))
	mutation := &api.Mutation{SetNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	return err
}

func (d *Dgraph) DeleteEdge(ctx context.Context, node1, edge, node2 string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	nquads := []byte(fmt.Sprintf("<%s> <%s> <%s> .", node1, edge, node2))
	mutation := &api.Mutation{DelNquads: nquads, CommitNow: true}
	_, err := d.client.NewTxn().Mutate(ctx, mutation)
	return err
}
