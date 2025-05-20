package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/utils/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEdge(c *gin.Context) {
	var request map[string][]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if len(request["source"]) != len(request["edge"]) ||
		len(request["source"]) != len(request["target"]) ||
		len(request["edge"]) != len(request["target"]) {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}
	if len(request["source"]) == 0 {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err = h.data.CreateEdge(ctx, request["source"], request["edge"], request["target"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) deleteEdge(c *gin.Context) {
	var request map[string][]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if len(request["source"]) != len(request["edge"]) ||
		len(request["source"]) != len(request["target"]) ||
		len(request["edge"]) != len(request["target"]) {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}
	if len(request["source"]) == 0 {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err = h.data.DeleteEdge(ctx, request["source"], request["edge"], request["target"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
