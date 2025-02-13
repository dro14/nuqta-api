package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	post.DType = []string{"Post"}
	post.Uid = "_:post"
	post.PostedAt = time.Now().Unix()

	assigned, err := d.setJson(ctx, post)
	if err != nil {
		return nil, err
	}

	post.Uid = assigned.Uids["post"]
	return post, nil
}

func (d *Dgraph) GetPost(ctx context.Context, uid, postUid string) (*models.Post, error) {
	query := fmt.Sprintf(postQuery, postUid)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	post := response["posts"][0]
	if post.PostedAt == 0 {
		return nil, e.ErrNotFound
	}

	post.Author, err = d.GetUserByUid(ctx, uid, post.Author.Uid)
	if err != nil {
		return nil, err
	}

	if post.InReplyTo != nil {
		post.InReplyTo.Author, err = d.GetUserByUid(ctx, uid, post.InReplyTo.Author.Uid)
		if err != nil {
			return nil, err
		}
	}

	post.IsReplied, err = d.IsReplied(ctx, uid, postUid)
	if err != nil {
		return nil, err
	}

	post.IsReposted, err = d.GetEdge(ctx, uid, "repost", postUid)
	if err != nil {
		return nil, err
	}

	post.IsLiked, err = d.GetEdge(ctx, uid, "like", postUid)
	if err != nil {
		return nil, err
	}

	post.IsClicked, err = d.GetEdge(ctx, uid, "click", postUid)
	if err != nil {
		return nil, err
	}

	post.IsViewed, err = d.GetEdge(ctx, uid, "view", postUid)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (d *Dgraph) GetPosts(ctx context.Context, uid string, postUids []string) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, len(postUids))
	for _, postUid := range postUids {
		post, err := d.GetPost(ctx, uid, postUid)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (d *Dgraph) GetRecentPosts(ctx context.Context) ([]*models.Post, error) {
	timestamp := time.Now().AddDate(0, 0, -2).Unix()
	query := fmt.Sprintf(recentPostsQuery, timestamp)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response["posts"], nil
}

func (d *Dgraph) GetFollowingPosts(ctx context.Context, uid string, before int64) ([]*models.Post, error) {
	query := fmt.Sprintf(followingQuery, uid, before)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response["posts"], nil
}

func (d *Dgraph) GetUserPosts(ctx context.Context, tab, userUid string, before int64) ([]string, error) {
	var query string
	switch tab {
	case "posts":
		query = fmt.Sprintf(userPostsQuery, userUid, before)
	case "replies":
		query = fmt.Sprintf(userRepliesQuery, userUid, before)
	case "reposts":
		query = fmt.Sprintf(userRepostsQuery, userUid, before)
	case "likes":
		query = fmt.Sprintf(userLikesQuery, userUid, before)
	}

	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
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

func (d *Dgraph) GetPopularReplies(ctx context.Context, postUid string, offset int) ([]string, error) {
	query := fmt.Sprintf(popularRepliesQuery, postUid)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	posts := response["posts"][0]["replies"]
	for i, post := range posts {
		posts[i].Score = 2.0*float64(post.Replies) +
			1.5*float64(post.Reposts) +
			1.0*float64(post.Likes) +
			0.5*float64(post.Clicks) +
			0.1*float64(post.Views)
	}

	slices.SortFunc(
		posts,
		func(a, b *models.Post) int {
			if a.Score < b.Score {
				return 1
			} else if a.Score > b.Score {
				return -1
			} else {
				return 0
			}
		},
	)

	var postUids []string
	for _, post := range posts[offset : offset+20] {
		postUids = append(postUids, post.Uid)
	}

	return postUids, nil
}

func (d *Dgraph) GetRecentReplies(ctx context.Context, postUid string, before int64) ([]string, error) {
	query := fmt.Sprintf(recentRepliesQuery, postUid, before)
	bytes, err := d.get(ctx, query)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
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

func (d *Dgraph) DeletePost(ctx context.Context, postUid string) error {
	query := fmt.Sprintf(`<%s> * * .`, postUid)
	return d.deleteNquads(ctx, query)
}
