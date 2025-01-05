package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getSchema(c *gin.Context) {
	schema, err := h.db.GetSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, schema)
}

func (h *Handler) updateSchema(c *gin.Context) {
	err := h.db.UpdateSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) deleteSchema(c *gin.Context) {
	err := h.db.DeleteSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Status(http.StatusOK)
}
