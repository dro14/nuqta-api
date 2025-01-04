package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	existingUser, err := h.db.GetUser(ctx, "firebase_uid", user.FirebaseUid)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusOK, existingUser)
		return
	}

	user, err = h.db.CreateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.search.AddUser(user)
	if err != nil {
		log.Printf("user %s: can't add user to search index: %s", user.Uid, err)
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := h.db.GetUser(ctx, c.Param("by"), c.Param("value"))
	if errors.Is(err, e.ErrUnknownParam) {
		c.Status(http.StatusBadRequest)
		return
	} else if errors.Is(err, e.ErrNotFound) {
		c.Status(http.StatusNotFound)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	err = h.db.UpdateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParam))
		return
	}

	ctx := c.Request.Context()
	err := h.db.DeleteUser(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.search.DeleteUser(uid)
	if err != nil {
		log.Printf("user %s: can't delete user from search index: %s", uid, err)
	}

	c.Status(http.StatusOK)
}
