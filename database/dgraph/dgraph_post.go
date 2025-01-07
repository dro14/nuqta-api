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

	nquads := fmt.Sprintf("<_:post> <author_uid> <%s> .", post.Author.Uid)
	post.Author = nil
	if post.InReplyTo != nil {
		nquads += fmt.Sprintf("\n<_:post> <in_reply_to_uid> <%s> .", post.InReplyTo.Uid)
		post.InReplyTo = nil
	}

	mutation := &api.Mutation{SetJson: json, SetNquads: []byte(nquads), CommitNow: true}
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
	for _, post := range response["posts"] {
		posts = append(posts, post.Uid)
	}

	return posts, nil
}

func (d *Dgraph) GetPost(ctx context.Context, uid string) (*models.Post, error) {
	query := fmt.Sprintf(postQuery, uid)
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
	query := fmt.Sprintf(userPostsQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var posts []string
	for _, post := range response["posts"] {
		posts = append(posts, post.Uid)
	}

	return posts, nil
}

func (d *Dgraph) GetPostReplies(ctx context.Context, uid string) ([]string, error) {
	query := fmt.Sprintf(postRepliesQuery, uid)
	resp, err := d.client.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(resp.Json, &response)
	if err != nil {
		return nil, err
	}

	var replies []string
	for _, reply := range response["replies"] {
		replies = append(replies, reply.Uid)
	}

	return replies, nil
}
