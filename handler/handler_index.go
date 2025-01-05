package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoQuery))
		return
	}

	users, err := h.index.SearchUser(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	ctx := c.Request.Context()
	for i := range users {
		user, err := h.db.GetUser(ctx, "uid", users[i].Uid)
		if err != nil {
			continue
		}
		users[i] = user
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) increment(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParam))
		return
	}

	err := h.index.IncrementUserHits(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}
