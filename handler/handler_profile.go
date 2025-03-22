package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

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

	if user.FirebaseUid != firebaseUid {
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

	username := strings.Split(user.Email, "@")[0]
	for {
		userUid, err := h.index.GetUidByUsername(username)
		if userUid != "" {
			lastCharIndex := len(username) - 1
			if lastCharIndex >= 0 && username[lastCharIndex] >= '0' && username[lastCharIndex] < '9' {
				username = username[:lastCharIndex] + string(username[lastCharIndex]+1)
			} else {
				username = username + "0"
			}
		} else if errors.Is(err, e.ErrNotFound) {
			user.Username = username
			break
		} else {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
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
	if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusNotFound, failure(err))
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		c.JSON(http.StatusOK, user)
	}
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

	if user.FirebaseUid != firebaseUid {
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
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" || request.Attribute == "" || request.Value == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	profile, err := h.db.GetProfile(ctx, firebaseUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if profile.Uid != request.Uid {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	err = h.db.DeleteProfileAttribute(ctx, request.Uid, request.Attribute, request.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if request.Attribute == "name" {
		err = h.index.DeleteName(request.Uid)
		if err != nil {
			log.Printf("user %s: can't delete name in search index: %s", request.Uid, err)
		}
	}
}

func (h *Handler) isAvailable(c *gin.Context) {
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
