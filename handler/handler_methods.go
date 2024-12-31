package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.Root)
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

func (h *Handler) PostUser(c *gin.Context) {
	pong, err := h.cache.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	schema, err := h.db.GetSchema()
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pong":   pong,
		"schema": schema,
	})
}

func (h *Handler) GetUser(c *gin.Context) {
	pong, err := h.cache.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	person, err := h.db.SetObject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pong":   pong,
		"person": person,
	})
}

func (h *Handler) PutUser(c *gin.Context) {

}

func (h *Handler) PatchUser(c *gin.Context) {

}

func (h *Handler) DeleteUser(c *gin.Context) {

}

func (h *Handler) SearchUser(c *gin.Context) {

}
