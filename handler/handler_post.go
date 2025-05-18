package handler

import (
	"net/http"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

type postRequest struct {
	Key      string   `json:"key"`
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
	request := &postRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	var postUids []string
	uid := c.GetString("uid")
	ctx := c.Request.Context()
	first, second, _ := strings.Cut(request.Key, ":")
	switch first {
	case "feed_following":
		posts, err := h.data.GetFollowingPosts(ctx, uid, request.Offset)
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
		postUids, err = h.data.GetReplies(ctx, uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "feed_saved":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetSavedPosts(ctx, uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "feed_history":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetViewedPosts(ctx, uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "user_posts", "user_replies", "user_media", "user_reposts", "user_likes":
		first = first[5:]
		if second == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetUserPosts(ctx, first, second, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_popular":
		if second == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetPopularReplies(ctx, second, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "replies_latest":
		if second == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err = h.data.GetLatestReplies(ctx, second, request.Before)
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
		post, err := h.data.GetPost(ctx, uid, postUid)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) editPost(c *gin.Context) {
	post := &models.Post{}
	err := c.ShouldBindJSON(post)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if post.Author.Uid != c.GetString("uid") {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	ctx := c.Request.Context()
	err = h.data.EditPost(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) deletePost(c *gin.Context) {
	var request map[string]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request["post_uid"] == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	uid := c.GetString("uid")
	ctx := c.Request.Context()
	_, err = h.data.GetPost(ctx, uid, request["post_uid"])
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	err = h.data.DeletePost(ctx, request["post_uid"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
