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
	"save",
	"report",
}

func (d *Dgraph) CreateEdge(ctx context.Context, source, edge, target []string) error {
	var objects []map[string]any
	for i := range source {
		if !slices.Contains(edges, edge[i]) {
			return e.ErrUnknownEdge
		}
		objects = append(objects, map[string]any{
			"uid": source[i],
			edge[i]: map[string]any{
				"uid":                  target[i],
				edge[i] + "|timestamp": time.Now().Unix(),
			},
		})
	}
	_, err := d.set(ctx, objects)
	return err
}

func (d *Dgraph) DeleteEdge(ctx context.Context, source, edge, target []string) error {
	if !slices.Contains(edges, edge[0]) {
		return e.ErrUnknownEdge
	}
	object := map[string]any{
		"uid": source[0],
		edge[0]: map[string]any{
			"uid": target[0],
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
