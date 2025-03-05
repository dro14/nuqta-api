package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUser(c *gin.Context) {
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

	if request.Username != "" {
		request.UserUid, err = h.index.GetUidByUsername(request.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	}

	if request.UserUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	user, err := h.db.GetUser(ctx, request.Uid, request.UserUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserList(c *gin.Context) {
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

	var userUids []string
	ctx := c.Request.Context()
	switch request.Tab {
	case "search":
		if request.Query == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.index.SearchUser(request.Query, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "followers":
		if request.UserUid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.db.GetUserFollowers(ctx, request.UserUid, request.After)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "following":
		if request.UserUid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.db.GetUserFollowing(ctx, request.UserUid, request.After)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	users := make([]*models.User, 0, len(userUids))
	for i := range userUids {
		user, err := h.db.GetUser(ctx, request.Uid, userUids[i])
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) hitUser(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.UserUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	err = h.index.HitUser(request.UserUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
