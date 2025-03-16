package dgraph

import (
	"context"
	"encoding/json"
	"slices"
	"strconv"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
)

func (d *Dgraph) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	post.DType = []string{"Post"}
	post.Uid = "_:post"
	post.PostedAt = time.Now().Unix()

	assigned, err := d.set(ctx, post)
	if err != nil {
		return nil, err
	}

	post.Uid = assigned.Uids["post"]
	return post, nil
}

func (d *Dgraph) GetPost(ctx context.Context, uid, postUid string, withInReplyTo bool) (*models.Post, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.get(ctx, postQuery, vars)
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

	vars = map[string]string{
		"$uid":      uid,
		"$post_uid": postUid,
	}
	bytes, err = d.get(ctx, postEdgesQuery, vars)
	if err != nil {
		return nil, err
	}

	var edges map[string][]map[string][]any
	err = json.Unmarshal(bytes, &edges)
	if err != nil {
		return nil, err
	}

	if len(edges["users"]) > 0 {
		user := edges["users"][0]
		post.IsReplied = len(user["is_replied"]) > 0
		post.IsReposted = len(user["is_reposted"]) > 0
		post.IsLiked = len(user["is_liked"]) > 0
		post.IsClicked = len(user["is_clicked"]) > 0
		post.IsViewed = len(user["is_viewed"]) > 0
		post.IsSaved = len(user["is_saved"]) > 0
	}

	if !withInReplyTo {
		post.InReplyTo = nil
	}

	return post, nil
}

func (d *Dgraph) GetPosts(ctx context.Context, uid string, postUids []string, withInReplyTo bool) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, len(postUids))
	for _, postUid := range postUids {
		post, err := d.GetPost(ctx, uid, postUid, withInReplyTo)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (d *Dgraph) GetLatestPosts(ctx context.Context) ([]*models.Post, error) {
	after := time.Now().AddDate(0, 0, -2).Unix()
	vars := map[string]string{
		"$after": strconv.FormatInt(after, 10),
	}
	bytes, err := d.get(ctx, latestPostsQuery, vars)
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
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.get(ctx, followingPostsQuery, vars)
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

func (d *Dgraph) GetReplies(ctx context.Context, uid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.get(ctx, repliesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var postUids []string
	for _, user := range response["users"] {
		for _, post := range user["posts"] {
			for _, reply := range post["replies"] {
				postUids = append(postUids, reply.Uid)
			}
		}
	}

	return postUids, nil
}

func (d *Dgraph) GetSavedPosts(ctx context.Context, uid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.get(ctx, savedPostsQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var postUids []string
	for _, user := range response["users"] {
		for _, post := range user["posts"] {
			postUids = append(postUids, post.Uid)
		}
	}

	return postUids, nil
}

func (d *Dgraph) GetUserPosts(ctx context.Context, tab, userUid string, before int64) ([]string, error) {
	var query string
	switch tab {
	case "posts":
		query = userPostsQuery
	case "replies":
		query = userRepliesQuery
	case "reposts":
		query = userRepostsQuery
	case "likes":
		query = userLikesQuery
	}

	vars := map[string]string{
		"$user_uid": userUid,
		"$before":   strconv.FormatInt(before, 10),
	}
	bytes, err := d.get(ctx, query, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var postUids []string
	for _, user := range response["users"] {
		for _, post := range user["posts"] {
			postUids = append(postUids, post.Uid)
		}
	}

	return postUids, nil
}

func (d *Dgraph) GetPopularReplies(ctx context.Context, postUid string, offset int64) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.get(ctx, postRepliesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	if len(response["posts"]) == 0 {
		return nil, nil
	}

	replies := response["posts"][0]["replies"]
	for i, reply := range replies {
		replies[i].Score = 2.0*float64(reply.Replies) +
			1.5*float64(reply.Reposts) +
			1.0*float64(reply.Likes) +
			0.5*float64(reply.Clicks) +
			0.1*float64(reply.Views) -
			1.0*float64(reply.Removes)
	}

	slices.SortFunc(
		replies,
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

	if offset > 0 {
		replies = replies[offset:]
	}

	if len(replies) > 20 {
		replies = replies[:20]
	}

	var replyUids []string
	for _, reply := range replies {
		replyUids = append(replyUids, reply.Uid)
	}

	return replyUids, nil
}

func (d *Dgraph) GetLatestReplies(ctx context.Context, postUid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
		"$before":   strconv.FormatInt(before, 10),
	}
	bytes, err := d.get(ctx, latestRepliesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var replyUids []string
	for _, post := range response["posts"] {
		for _, reply := range post["replies"] {
			replyUids = append(replyUids, reply.Uid)
		}
	}

	return replyUids, nil
}

func (d *Dgraph) GetPostReplies(ctx context.Context, postUid string) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.get(ctx, postRepliesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var replyUids []string
	for _, post := range response["posts"] {
		for _, reply := range post["replies"] {
			replyUids = append(replyUids, reply.Uid)
		}
	}

	return replyUids, nil
}

func (d *Dgraph) DeletePost(ctx context.Context, postUid string) error {
	replyUids, err := d.GetPostReplies(ctx, postUid)
	if err != nil {
		return err
	}

	var objects []map[string]any
	for _, replyUid := range replyUids {
		objects = append(objects, map[string]any{
			"uid": replyUid,
			"in_reply_to": map[string]any{
				"uid": postUid,
			},
		})
	}

	objects = append(objects, map[string]any{
		"uid": postUid,
	})
	return d.delete(ctx, objects)
}
