package data

import (
	"context"
	"slices"
	"time"

	"github.com/dro14/nuqta-service/e"
)

var edges = []string{
	"follow",
	"block",
	"chat",
	"repost",
	"like",
	"click",
	"view",
	"save",
}

func (d *Data) CreateEdge(ctx context.Context, source, edge, target []string) error {
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
	_, err := d.graphSet(ctx, objects)
	return err
}

func (d *Data) DeleteEdge(ctx context.Context, source, edge, target []string) error {
	var objects []map[string]any
	for i := range source {
		if !slices.Contains(edges, edge[i]) {
			return e.ErrUnknownEdge
		}
		objects = append(objects, map[string]any{
			"uid": source[i],
			edge[i]: map[string]any{
				"uid": target[i],
			},
		})
	}
	return d.graphDelete(ctx, objects)
}
