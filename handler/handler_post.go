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
}

func (h *Handler) getAllPosts(c *gin.Context) {
	ctx := c.Request.Context()
	allPosts, err := h.db.GetAllPosts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, allPosts)
}

func (h *Handler) getForYouPosts(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	postUids := h.rec.GetRecs()

	ctx := c.Request.Context()
	posts := make([]*models.Post, 0, 20)
	for _, uid := range postUids {
		isViewed, err := h.db.DoesEdgeExist(ctx, uid, "~view", firebaseUid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		} else if isViewed {
			continue
		}
		post, err := h.db.GetPostByUid(ctx, firebaseUid, uid)
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
}

func (h *Handler) getPost(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.db.GetPostByUid(ctx, firebaseUid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) getFollowingPosts(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	before := c.Param("before")
	if before == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	postUids, err := h.db.GetFollowingPosts(ctx, firebaseUid, before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	posts := make([]*models.Post, 0, len(postUids))
	for _, uid := range postUids {
		post, err := h.db.GetPostByUid(ctx, firebaseUid, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
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

func (h *Handler) deletePost(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.db.GetPostByUid(ctx, firebaseUid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	author, err := h.db.GetUserByUid(ctx, firebaseUid, post.Author.Uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if author.FirebaseUid != firebaseUid {
		c.JSON(http.StatusForbidden, failure(e.ErrForbidden))
		return
	}

	err = h.db.DeletePost(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
