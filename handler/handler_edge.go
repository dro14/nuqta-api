package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEdge(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" || request.Edge == "" || request.Target == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err = h.db.CreateEdge(ctx, request.Uid, request.Edge, request.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) deleteEdge(c *gin.Context) {
	request := &models.Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.Uid == "" || request.Edge == "" || request.Target == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err = h.db.DeleteEdge(ctx, request.Uid, request.Edge, request.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
