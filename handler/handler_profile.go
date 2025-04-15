package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

const version = "1.0.4"

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
	existingUser, err := h.data.GetProfile(ctx, firebaseUid)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if existingUser != nil {
		existingUser.Version = version
		c.JSON(http.StatusOK, existingUser)
		return
	}

	if user.InvitedBy != nil {
		if user.InvitedBy.Uid != "" {
			if strings.HasPrefix(user.InvitedBy.Uid, "0x") {
				_, err = strconv.ParseInt(user.InvitedBy.Uid[2:], 16, 64)
				if err != nil {
					user.InvitedBy = nil
				}
			} else {
				osVersion := user.InvitedBy.Uid
				userUid, err := h.data.GetReferrer(ctx, c.ClientIP(), osVersion)
				if userUid != "" {
					user.InvitedBy = &models.User{Uid: userUid}
				} else {
					if err != nil {
						log.Printf("can't get referrer: %s", err)
					}
					user.InvitedBy = nil
				}
			}
		} else {
			user.InvitedBy = nil
		}
	}

	err = h.data.CreateProfile(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	user.Version = version
	c.JSON(http.StatusOK, user)
}

func (h *Handler) getProfile(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	user, err := h.data.GetProfile(ctx, firebaseUid)
	if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusNotFound, failure(err))
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		user.Version = version
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
	err = h.data.UpdateProfile(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) isAvailable(c *gin.Context) {
	var request map[string]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request["username"] == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	uid, err := h.data.GetUidByUsername(ctx, request["username"])
	if uid != "" {
		c.JSON(http.StatusOK, gin.H{"available": false})
	} else if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusOK, gin.H{"available": true})
	} else {
		c.JSON(http.StatusInternalServerError, failure(err))
	}
}
