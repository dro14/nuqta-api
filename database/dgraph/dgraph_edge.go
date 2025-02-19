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
	"save",
}

func (d *Dgraph) CreateEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	now := time.Now().Unix()
	query := fmt.Sprintf("<%s> <%s> <%s> (timestamp=%d) .", source, edge, target, now)
	return d.setNquads(ctx, query)
}

func (d *Dgraph) DeleteEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	query := fmt.Sprintf("<%s> <%s> <%s> .", source, edge, target)
	return d.deleteNquads(ctx, query)
}

func (d *Dgraph) IsPostViewed(ctx context.Context, uid, postUid string) (bool, error) {
	vars := map[string]string{
		"$uid":      uid,
		"$post_uid": postUid,
	}
	bytes, err := d.get(ctx, isViewedQuery, vars)
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
