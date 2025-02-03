package handler

import (
	"errors"
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUser(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var uid string
	if c.Query("username") != "" {
		var err error
		uid, err = h.index.GetUidByUsername(c.Query("username"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	} else {
		uid = c.Query("uid")
		if uid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
	}

	ctx := c.Request.Context()
	user, err := h.db.GetUserByUid(ctx, firebaseUid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) isUsernameAvailable(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	uid, err := h.index.GetUidByUsername(username)
	if uid != "" {
		c.JSON(http.StatusOK, gin.H{"available": false})
	} else if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusOK, gin.H{"available": true})
	} else {
		c.JSON(http.StatusInternalServerError, failure(err))
	}
}
