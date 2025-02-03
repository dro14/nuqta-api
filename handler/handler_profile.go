package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createProfile(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if firebaseUid != user.FirebaseUid {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	existingUser, err := h.db.GetProfile(ctx, firebaseUid)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusOK, existingUser)
		return
	}

	user, err = h.db.CreateProfile(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.index.AddUser(user)
	if err != nil {
		log.Printf("user %s: can't add user to search index: %s", user.Uid, err)
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getProfile(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	user, err := h.db.GetProfile(ctx, firebaseUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateProfile(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if firebaseUid != user.FirebaseUid {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	err = h.db.UpdateProfile(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.index.UpdateUser(user)
	if err != nil {
		log.Printf("user %s: can't update user in search index: %s", user.Uid, err)
	}
}

func (h *Handler) deleteProfileAttribute(c *gin.Context) {
	uid := c.Param("uid")
	attribute := c.Param("attribute")
	value := c.Query("value")
	if uid == "" || attribute == "" || value == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.DeleteProfileAttribute(ctx, uid, attribute, value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if attribute == "name" {
		err = h.index.DeleteName(uid)
		if err != nil {
			log.Printf("user %s: can't delete name in search index: %s", uid, err)
		}
	}
}
