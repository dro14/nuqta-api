package handler

import (
	"errors"
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

	if request.Uid == "" && request.Username == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	userUid, err := h.index.GetUidByUsername(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	ctx := c.Request.Context()
	user, err := h.db.GetUserByUid(ctx, request.Uid, userUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) isUsernameAvailable(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	uid, err := h.index.GetUidByUsername(request.Username)
	if uid != "" {
		c.JSON(http.StatusOK, gin.H{"available": false})
	} else if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusOK, gin.H{"available": true})
	} else {
		c.JSON(http.StatusInternalServerError, failure(err))
	}
}

func (h *Handler) searchUser(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" || request.Query == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	userUids, err := h.index.SearchUser(request.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	ctx := c.Request.Context()
	users := make([]*models.User, 0, len(userUids))
	for i := range userUids {
		user, err := h.db.GetUserByUid(ctx, request.Uid, userUids[i])
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) updateUser(c *gin.Context) {
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

	err = h.index.IncrementHits(request.UserUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
