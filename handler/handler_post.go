package handler

import (
	"net/http"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

type PostRequest struct {
	Uid      string   `json:"uid"`
	Tab      string   `json:"tab"`
	UserUid  string   `json:"user_uid"`
	PostUid  string   `json:"post_uid"`
	PostUids []string `json:"post_uids"`
	Before   int64    `json:"before"`
	Offset   int64    `json:"offset"`
}

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
	request := &PostRequest{}
	err := c.ShouldBindJSON(request)
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
	case "feed_following":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		posts, err := h.data.GetFollowingPosts(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		c.JSON(http.StatusOK, posts)
		return
	case "feed_replies":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetReplies(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "feed_saved":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetSavedPosts(ctx, request.Uid, request.Before)
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
		postUids, err = h.data.GetUserPosts(ctx, request.Tab, request.UserUid, request.Before)
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
		postUids, err = h.data.GetPopularReplies(ctx, request.PostUid, request.Offset)
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
		postUids, err = h.data.GetLatestReplies(ctx, request.PostUid, request.Before)
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

	posts := make([]*models.Post, 0, len(postUids))
	for _, postUid := range postUids {
		post, err := h.data.GetPost(ctx, request.Uid, postUid, withInReplyTo)
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
