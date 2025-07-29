package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getSchema(c *gin.Context) {
	schema, err := h.data.GetSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, schema)
}

func (h *Handler) updateSchema(c *gin.Context) {
	err := h.data.UpdateSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) deleteSchema(c *gin.Context) {
	err := h.data.DeleteSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
