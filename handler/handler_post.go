package handler

import (
	"net/http"
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
	err = h.data.CreatePost(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) getPostList(c *gin.Context) {
	var request map[string]any
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	uid, ok := request["uid"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var postUids []string
	withInReplyTo := true
	ctx := c.Request.Context()
	switch request["tab"] {
	case "feed_following":
		before, ok := request["before"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		posts, err := h.data.GetFollowingPosts(ctx, uid, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		c.JSON(http.StatusOK, posts)
		return
	case "feed_replies":
		before, ok := request["before"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetReplies(ctx, uid, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "feed_saved":
		before, ok := request["before"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetSavedPosts(ctx, uid, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "user_posts", "user_replies", "user_reposts", "user_likes":
		tab, ok := request["tab"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		tab = strings.TrimPrefix(tab, "user_")
		userUid, ok := request["user_uid"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		before, ok := request["before"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetUserPosts(ctx, tab, userUid, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_popular":
		withInReplyTo = false
		postUid, ok := request["post_uid"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		offset, ok := request["offset"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetPopularReplies(ctx, postUid, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_latest":
		withInReplyTo = false
		postUid, ok := request["post_uid"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		before, ok := request["before"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetLatestReplies(ctx, postUid, before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		list, ok := request["post_uids"].([]any)
		if ok && len(list) > 0 {
			for _, postUid := range list {
				postUids = append(postUids, postUid.(string))
			}
		} else {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
	}

	posts := make([]*models.Post, 0, len(postUids))
	for _, postUid := range postUids {
		post, err := h.data.GetPost(ctx, uid, postUid, withInReplyTo)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) deletePost(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var request map[string]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request["uid"] == "" || request["post_uid"] == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.data.GetPost(ctx, request["uid"], request["post_uid"], false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if request["uid"] != post.Author.Uid {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	author, err := h.data.GetProfile(ctx, firebaseUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if request["uid"] != author.Uid {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	err = h.data.DeletePost(ctx, request["post_uid"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
