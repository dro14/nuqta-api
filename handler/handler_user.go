package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Uid      string   `json:"uid"`
	Tab      string   `json:"tab"`
	Query    string   `json:"query"`
	UserUid  string   `json:"user_uid"`
	UserUids []string `json:"user_uids"`
	After    string   `json:"after"`
	Offset   int64    `json:"offset"`
}

func (h *Handler) getUserList(c *gin.Context) {
	request := &UserRequest{}
	err := c.ShouldBindJSON(request)
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
			c.JSON(http.StatusOK, make([]*models.User, 0))
			return
		}
		userUids, err = h.data.SearchUser(ctx, request.Query, request.Offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "followers", "following":
		if request.UserUid == "" {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		reverse := request.Tab == "followers"
		userUids, err = h.data.GetUserFollows(ctx, request.UserUid, request.After, reverse)
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
		user, err := h.data.GetUser(ctx, request.Uid, userUid)
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

	if request["uid"] == "" || request["username"] == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	userUid, err := h.data.GetUidByUsername(ctx, request["username"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	user, err := h.data.GetUser(ctx, request["uid"], userUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}
