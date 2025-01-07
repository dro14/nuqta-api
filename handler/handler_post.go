package handler

import (
	"net/http"

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
	_, err = h.db.CreatePost(ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getPosts(c *gin.Context) {
	ctx := c.Request.Context()
	posts, err := h.db.GetPosts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) getPost(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.db.GetPost(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	firebaseUid, ok := ctx.Value("firebase_uid").(string)
	if ok {
		post.IsLiked, err = h.db.DoesEdgeExist(ctx, firebaseUid, "like", uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}

		post.IsReposted, err = h.db.DoesEdgeExist(ctx, firebaseUid, "repost", uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) getUserPosts(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	posts, err := h.db.GetUserPosts(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) getPostReplies(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	replies, err := h.db.GetPostReplies(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, replies)
}
