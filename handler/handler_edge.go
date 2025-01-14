package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEdge(c *gin.Context) {
	source := c.Param("source")
	edge := c.Param("edge")
	target := c.Param("target")
	if source == "" || edge == "" || target == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.CreateEdge(ctx, source, edge, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) deleteEdge(c *gin.Context) {
	source := c.Param("source")
	edge := c.Param("edge")
	target := c.Param("target")
	if source == "" || edge == "" || target == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.DeleteEdge(ctx, source, edge, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
