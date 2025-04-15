package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserList(c *gin.Context) {
	var request map[string]any
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	uid, ok := request["uid"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var userUids []string
	ctx := c.Request.Context()
	switch request["tab"] {
	case "search":
		query, ok := request["query"].(string)
		if !ok {
			c.JSON(http.StatusOK, make([]*models.User, 0))
			return
		}
		offset, ok := request["offset"].(int64)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		userUids, err = h.data.SearchUser(ctx, query, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	case "followers", "following":
		userUid, ok := request["user_uid"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		after, ok := request["after"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
			return
		}
		reverse := request["tab"] == "followers"
		userUids, err = h.data.GetUserFollows(ctx, userUid, after, reverse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
	default:
		list, ok := request["user_uids"].([]any)
		if ok && len(list) > 0 {
			for _, userUid := range list {
				userUids = append(userUids, userUid.(string))
			}
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
