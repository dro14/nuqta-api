package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/dro14/nuqta-service/e"
)

var edges = []string{
	"follow",
	"repost",
	"like",
	"click",
	"view",
	"remove",
}

func (d *Dgraph) CreateEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	now := time.Now().Unix()
	query := fmt.Sprintf("<%s> <%s> <%s> (timestamp=%d) .", source, edge, target, now)
	return d.setNquads(ctx, query)
}

func (d *Dgraph) GetEdge(ctx context.Context, source, edge, target string) (bool, error) {
	if !slices.Contains(edges, edge) {
		return false, e.ErrUnknownEdge
	}

	query := fmt.Sprintf(edgeQuery, source, edge, target)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return false, err
	}

	var response map[string][]any
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return false, err
	}

	return len(response["edges"]) > 0, nil
}

func (d *Dgraph) IsReplied(ctx context.Context, userUid, postUid string) (bool, error) {
	query := fmt.Sprintf(isRepliedQuery, userUid, postUid)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return false, err
	}

	var response map[string][]any
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return false, err
	}

	return len(response["edges"]) > 0, nil
}

func (d *Dgraph) DeleteEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	query := fmt.Sprintf("<%s> <%s> <%s> .", source, edge, target)
	return d.deleteNquads(ctx, query)
}
