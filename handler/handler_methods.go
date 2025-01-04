package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.Root)

	h.engine.GET("/schema", h.GetSchema)
	h.engine.PUT("/schema", h.SetSchema)
	h.engine.DELETE("/schema", h.DeleteSchema)

	h.engine.POST("/user", h.CreateUser)
	h.engine.GET("/user/:by/:value", h.GetUser)
	h.engine.PUT("/user", h.UpdateUser)
	h.engine.DELETE("/user/:uid", h.DeleteUser)

	h.engine.GET("/search", h.SearchUser)
	h.engine.PATCH("/increment_hits/:uid", h.IncrementHits)
	return h.engine.Run(":" + port)
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (h *Handler) GetSchema(c *gin.Context) {
	schema, err := h.db.GetSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.String(http.StatusOK, schema)
}

func (h *Handler) SetSchema(c *gin.Context) {
	err := h.db.UpdateSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) DeleteSchema(c *gin.Context) {
	err := h.db.DeleteSchema(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.Status(http.StatusOK)
}
