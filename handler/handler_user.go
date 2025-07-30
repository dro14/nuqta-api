package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dro14/nuqta-api/models"
	"github.com/dro14/nuqta-api/utils/e"
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
	uid := c.GetString("uid")
	ctx := c.Request.Context()
	first, second, _ := strings.Cut(request.Key, ":")
	switch first {
	case "search":
		if second == "" {
			userUids, err = h.data.GetUserRecommendations(ctx, uid, request.Offset)
		} else {
			userUids, err = h.data.SearchUser(ctx, second, request.Offset)
		}
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
	case "invitations":
		userUids, err = h.data.GetUserInvitations(ctx, uid, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "post_reposts", "post_likes":
		first = first[5 : len(first)-1]
		if second == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.data.GetPostUsers(ctx, first, second, request.Offset)
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
		user, err := h.data.GetUser(ctx, uid, userUid)
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
	if errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusOK, &models.User{})
		return
	} else if err != nil {
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
