package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createEdge(c *gin.Context) {
	node1 := c.Param("node1")
	edge := c.Param("edge")
	node2 := c.Param("node2")
	if node1 == "" || edge == "" || node2 == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.CreateEdge(ctx, node1, edge, node2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) deleteEdge(c *gin.Context) {
	node1 := c.Param("node1")
	edge := c.Param("edge")
	node2 := c.Param("node2")
	if node1 == "" || edge == "" || node2 == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.DeleteEdge(ctx, node1, edge, node2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}
