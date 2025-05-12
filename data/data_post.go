package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"slices"
	"strconv"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/lib/pq"
)

func (d *Data) CreatePost(ctx context.Context, post *models.Post) error {
	post.Timestamp = time.Now().UnixMilli()

	object := map[string]any{
		"dgraph.type": "post",
		"uid":         "_:post",
		"timestamp":   post.Timestamp,
		"author": map[string]string{
			"uid": post.Author.Uid,
		},
	}

	if len(post.Images) > 0 {
		object["has_media"] = true
	}

	if post.InReplyTo != nil {
		object["in_reply_to"] = map[string]string{
			"uid": post.InReplyTo.Uid,
		}
	}

	assigned, err := d.graphSet(ctx, object)
	if err != nil {
		return err
	}
	post.Uid = assigned.Uids["post"]

	if len(post.Images) > 5 {
		post.Images = post.Images[:5]
	}

	err = d.dbExec(ctx,
		"INSERT INTO posts (id, timestamp, text, who_can_reply, images) VALUES ($1, $2, $3, $4, $5)",
		post.Uid, post.Timestamp, post.Text, post.WhoCanReply, pq.Array(post.Images),
	)
	if err != nil {
		object = map[string]any{"uid": post.Uid}
		d.graphDelete(ctx, object)
		return err
	}

	return nil
}

func (d *Data) GetPost(ctx context.Context, uid, postUid string) (*models.Post, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.graphGet(ctx, postQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	post := response["posts"][0]
	post.Uid = postUid

	var edited sql.NullInt64

	err = d.dbQueryRow(ctx,
		"SELECT timestamp, text, who_can_reply, images, edited FROM posts WHERE id = $1",
		[]any{postUid},
		&post.Timestamp, &post.Text, &post.WhoCanReply, pq.Array(&post.Images), &edited,
	)
	if err != nil {
		return nil, err
	}

	if edited.Valid {
		post.Edited = edited.Int64
	}

	vars = map[string]string{
		"$uid":      uid,
		"$post_uid": postUid,
	}
	bytes, err = d.graphGet(ctx, postEdgesQuery, vars)
	if err != nil {
		return nil, err
	}

	var edges map[string][]map[string][]map[string]string
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

	return post, nil
}

func (d *Data) GetFollowingPosts(ctx context.Context, uid string, offset int64) ([]*models.Post, error) {
	vars := map[string]string{
		"$uid": uid,
	}
	bytes, err := d.graphGet(ctx, followingQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]map[string][]*models.Following
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for _, user := range response["users"] {
		for _, following := range user["following"] {
			posts = append(posts, following.Posts...)
			for _, repost := range following.Reposts {
				posts = append(posts, &models.Post{
					Uid:       repost.Uid,
					Timestamp: repost.RepostsTimestamp,
					RepostedBy: &models.User{
						Uid: following.Uid,
					},
				})
			}
		}
	}

	slices.SortFunc(
		posts,
		func(a, b *models.Post) int {
			if a.Timestamp < b.Timestamp {
				return 1
			} else if a.Timestamp > b.Timestamp {
				return -1
			} else {
				return 0
			}
		},
	)

	var added []string
	var unique []*models.Post
	for _, post := range posts {
		if slices.Contains(added, post.Uid) {
			continue
		}
		added = append(added, post.Uid)
		unique = append(unique, post)
	}

	if offset > 0 {
		unique = unique[offset:]
	}

	result := make([]*models.Post, 0)
	for i := range unique {
		post, err := d.GetPost(ctx, uid, unique[i].Uid)
		if err != nil {
			return nil, err
		}
		if unique[i].RepostedBy != nil {
			post.RepostedBy = unique[i].RepostedBy
		}
		result = append(result, post)
		if len(result) == 20 {
			break
		}
	}

	return result, nil
}

func (d *Data) GetReplies(ctx context.Context, uid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.graphGet(ctx, repliesQuery, vars)
	if err != nil {
		return nil, err
	}

	var response map[string][]*models.Post
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	var postUids []string
	for _, reply := range response["replies"] {
		postUids = append(postUids, reply.Uid)
	}

	return postUids, nil
}

func (d *Data) GetSavedPosts(ctx context.Context, uid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.graphGet(ctx, savedQuery, vars)
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

func (d *Data) GetViewedPosts(ctx context.Context, uid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$uid":    uid,
		"$before": strconv.FormatInt(before, 10),
	}
	bytes, err := d.graphGet(ctx, viewedQuery, vars)
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

func (d *Data) GetUserPosts(ctx context.Context, tab, userUid string, before int64) ([]string, error) {
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
	bytes, err := d.graphGet(ctx, query, vars)
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

func (d *Data) GetPopularReplies(ctx context.Context, postUid string, offset int64) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.graphGet(ctx, postRepliesQuery, vars)
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
			1.0*float64(reply.Reports)
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

func (d *Data) GetLatestReplies(ctx context.Context, postUid string, before int64) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
		"$before":   strconv.FormatInt(before, 10),
	}
	bytes, err := d.graphGet(ctx, latestRepliesQuery, vars)
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

func (d *Data) GetPostReplies(ctx context.Context, postUid string) ([]string, error) {
	vars := map[string]string{
		"$post_uid": postUid,
	}
	bytes, err := d.graphGet(ctx, postRepliesQuery, vars)
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

func (d *Data) EditPost(ctx context.Context, post *models.Post) error {
	post.Edited = time.Now().UnixMilli()
	return d.dbExec(ctx,
		"UPDATE posts SET text = $1, who_can_reply = $2, images = $3, edited = $4 WHERE id = $5",
		post.Text, post.WhoCanReply, pq.Array(post.Images), post.Edited, post.Uid,
	)
}

func (d *Data) DeletePost(ctx context.Context, postUid string) error {
	replyUids, err := d.GetPostReplies(ctx, postUid)
	if err != nil {
		return err
	}
	for _, replyUid := range replyUids {
		err = d.DeletePost(ctx, replyUid)
		if err != nil {
			return err
		}
	}
	object := map[string]any{"uid": postUid}
	err = d.graphDelete(ctx, object)
	if err != nil {
		return err
	}
	return d.dbExec(ctx, "DELETE FROM posts WHERE id = $1", postUid)
}
