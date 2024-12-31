package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.Root)

	h.engine.GET("/schema", h.GetSchema)
	h.engine.PUT("/schema", h.SetSchema)

	h.engine.POST("/user", h.PostUser)
	h.engine.GET("/user", h.GetUser)
	h.engine.PUT("/user", h.PutUser)
	h.engine.PATCH("/user", h.PatchUser)
	h.engine.DELETE("/user", h.DeleteUser)

	h.engine.GET("/search", h.SearchUser)
	return h.engine.Run(":" + port)
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (h *Handler) GetSchema(c *gin.Context) {
	ctx := c.Request.Context()
	schema, err := h.db.ReadSchema(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, schema)
}

func (h *Handler) SetSchema(c *gin.Context) {
	ctx := c.Request.Context()
	err := h.db.UpdateSchema(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) PostUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	err = h.db.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *Handler) GetUser(c *gin.Context) {}

func (h *Handler) PutUser(c *gin.Context) {}

func (h *Handler) PatchUser(c *gin.Context) {}

func (h *Handler) DeleteUser(c *gin.Context) {}

func (h *Handler) SearchUser(c *gin.Context) {}
