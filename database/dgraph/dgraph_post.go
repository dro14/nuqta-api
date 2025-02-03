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

func (d *Dgraph) GetPosts(ctx context.Context) ([]string, error) {
	resp, err := d.client.NewTxn().Query(ctx, postsQuery)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var posts []string
	for _, post := range response["all_posts"] {
		posts = append(posts, post.Uid)
	}

	return posts, nil
}

func (d *Dgraph) GetPostByUid(ctx context.Context, firebaseUid, uid string) (*models.Post, error) {
	query := fmt.Sprintf(postByUidQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	post := response["posts"][0]
	if post.PostedAt == 0 {
		return nil, e.ErrNotFound
	}

	post.IsLiked, err = d.doesEdgeExist(ctx, firebaseUid, "like", uid)
	if err != nil {
		return nil, err
	}

	post.IsReposted, err = d.doesEdgeExist(ctx, firebaseUid, "repost", uid)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (d *Dgraph) GetUserPosts(ctx context.Context, uid string) ([]string, error) {
	query := fmt.Sprintf(userPostsQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var posts []string
	for _, user := range response["users"] {
		for _, post := range user["posts"] {
			posts = append(posts, post.Uid)
		}
	}

	return posts, nil
}

func (d *Dgraph) GetPostReplies(ctx context.Context, uid string) ([]string, error) {
	query := fmt.Sprintf(postRepliesQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var replies []string
	for _, post := range response["posts"] {
		for _, reply := range post["replies"] {
			replies = append(replies, reply.Uid)
		}
	}

	return replies, nil
}
