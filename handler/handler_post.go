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

func (h *Handler) getPost(c *gin.Context) {
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

	ctx := c.Request.Context()
	var posts []*models.Post
	switch request.Tab {
	case "feed_for_you":
		posts = make([]*models.Post, 0, 20)
		postUids := h.rec.GetRecs()
		for _, postUid := range postUids {
			isViewed, err := h.db.GetEdge(ctx, request.Uid, "view", postUid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			} else if isViewed {
				continue
			}
			post, err := h.db.GetPost(ctx, request.Uid, postUid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			posts = append(posts, post)
			if len(posts) == cap(posts) {
				break
			}
		}
	case "feed_following":
		if request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		posts, err = h.db.GetFollowingPosts(ctx, request.Uid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		for i, post := range posts {
			posts[i], err = h.db.GetPost(ctx, request.Uid, post.Uid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			posts[i].RepostedBy = post.RepostedBy
		}
	case "user_posts", "user_replies", "user_reposts", "user_likes":
		if request.UserUid == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err := h.db.GetUserPosts(ctx, request.Tab, request.UserUid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		posts, err = h.db.GetPosts(ctx, request.Uid, postUids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "reply_popular":
		if request.PostUid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err := h.db.GetPopularReplies(ctx, request.PostUid, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		posts, err = h.db.GetPosts(ctx, request.Uid, postUids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "reply_recent":
		if request.PostUid == "" || request.Before == 0 {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		postUids, err := h.db.GetRecentReplies(ctx, request.PostUid, request.Before)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		posts, err = h.db.GetPosts(ctx, request.Uid, postUids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		if len(request.PostUids) > 0 {
			posts, err = h.db.GetPosts(ctx, request.Uid, request.PostUids)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) deletePost(c *gin.Context) {
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

	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	post, err := h.db.GetPost(ctx, request.Uid, request.PostUid)
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
