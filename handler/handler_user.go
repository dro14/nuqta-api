package handler

import (
	"net/http"
	"strings"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Key      string   `json:"key"`
	UserUids []string `json:"user_uids"`
	Offset   int64    `json:"offset"`
}

func (h *Handler) getUserList(c *gin.Context) {
	request := &userRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	var userUids []string
	ctx := c.Request.Context()
	first, second, _ := strings.Cut(request.Key, ":")
	switch first {
	case "search":
		if second == "" {
			c.JSON(http.StatusOK, make([]*models.User, 0))
			return
		}
		userUids, err = h.data.SearchUser(ctx, second, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "followers", "following":
		if second == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.data.GetUserFollows(ctx, second, request.Offset, first == "followers")
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "invites":
		userUids, err = h.data.GetUserInvites(ctx, c.GetString("uid"), request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		if len(request.UserUids) > 0 {
			userUids = request.UserUids
		} else {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
	}

	users := make([]*models.User, 0, len(userUids))
	for _, userUid := range userUids {
		user, err := h.data.GetUser(ctx, c.GetString("uid"), userUid)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUserByUsername(c *gin.Context) {
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
	userUid, err := h.data.GetUidByUsername(ctx, request["username"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	user, err := h.data.GetUser(ctx, c.GetString("uid"), userUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
