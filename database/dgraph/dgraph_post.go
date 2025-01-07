package dgraph

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	post.DType = []string{"Post"}
	post.Uid = "_:post"
	json, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}

	mutation := &api.Mutation{SetJson: json, CommitNow: true}
	assigned, err := d.client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return nil, err
	}

	post.Uid = assigned.Uids["post"]
	return post, nil
}

func (d *Dgraph) GetPost(ctx context.Context, uid string) (*models.Post, error) {
	query := fmt.Sprintf(postsQuery, fmt.Sprintf(functions["uid"], uid))
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	if len(response["posts"]) > 0 && response["posts"][0].PostedAt != 0 {
		return response["posts"][0], nil
	} else {
		return nil, e.ErrNotFound
	}
}

func (d *Dgraph) GetUserPosts(ctx context.Context, uid string) ([]string, error) {
	query := fmt.Sprintf(postsOfUserQuery, fmt.Sprintf(functions["uid"], uid))
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]map[string]string
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var posts []string
	for _, post := range response["users"][0]["posts"] {
		posts = append(posts, post["uid"])
	}

	return posts, nil
}
