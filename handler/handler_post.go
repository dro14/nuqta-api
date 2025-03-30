package handler

import (
	"net/http"
	"slices"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createPost(c *gin.Context) {
	post := &models.Post{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	post, err = h.db.CreatePost(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) getPostList(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var postUids []string
	withInReplyTo := true
	ctx := c.Request.Context()
	switch request.Tab {
	case "feed_for_you":
		postUids = h.rec.GetRecs()
		posts := make([]*models.Post, 0, 20)
		for _, postUid := range postUids {
			post, err := h.db.GetPost(ctx, request.Uid, postUid, withInReplyTo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			posts = append(posts, post)
			if len(posts) == 20 {
				break
			}
		}
		c.JSON(http.StatusOK, posts)
		return
	case "feed_following":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		posts, err := h.db.GetFollowingPosts(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		for i, post := range posts {
			if len(post.Reposted) > 0 {
				posts[i].Timestamp = post.Reposted[0].RepostedTimestamp
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
		if len(posts) > 20 {
			posts = posts[:20]
		}
		for i, post := range posts {
			posts[i], err = h.db.GetPost(ctx, request.Uid, post.Uid, withInReplyTo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			if len(post.Reposted) > 0 {
				posts[i].RepostedBy = &models.User{Uid: post.Reposted[0].Uid}
			}
		}
		c.JSON(http.StatusOK, posts)
		return
	case "feed_replies":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.db.GetReplies(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "feed_saved":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.db.GetSavedPosts(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "user_posts", "user_replies", "user_reposts", "user_likes":
		request.Tab = strings.TrimPrefix(request.Tab, "user_")
		if request.UserUid == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.db.GetUserPosts(ctx, request.Tab, request.UserUid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_popular":
		withInReplyTo = false
		if request.PostUid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.db.GetPopularReplies(ctx, request.PostUid, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_latest":
		withInReplyTo = false
		if request.PostUid == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.db.GetLatestReplies(ctx, request.PostUid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		if len(request.PostUids) > 0 {
			postUids = request.PostUids
		} else {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
	}

	posts, err := h.db.GetPosts(ctx, request.Uid, postUids, withInReplyTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) deletePost(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" || request.PostUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.db.GetPost(ctx, request.Uid, request.PostUid, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if request.Uid != post.Author.Uid {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	author, err := h.db.GetProfile(ctx, firebaseUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if request.Uid != author.Uid {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	err = h.db.DeletePost(ctx, request.PostUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
