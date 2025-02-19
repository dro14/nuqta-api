package dgraph

import (
	"context"
	"encoding/json"
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
	object := map[string]any{
		"uid": source,
		edge: map[string]any{
			"uid":               target,
			edge + "|timestamp": time.Now().Unix(),
		},
	}
	_, err := d.set(ctx, object)
	return err
}

func (d *Dgraph) DeleteEdge(ctx context.Context, source, edge, target string) error {
	if !slices.Contains(edges, edge) {
		return e.ErrUnknownEdge
	}
	object := map[string]any{
		"uid": source,
		edge: map[string]any{
			"uid": target,
		},
	}
	return d.delete(ctx, object)
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

	return len(response["is_viewed"]) > 0, nil
}
