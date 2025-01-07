package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.root)

	group := h.engine.Group("/schema")
	group.GET("", h.getSchema)
	group.PUT("", h.updateSchema)
	group.DELETE("", h.deleteSchema)

	group = h.engine.Group("/user")
	group.POST("", h.createUser)
	group.GET("/:by/:value", h.getUser)
	group.PUT("", h.updateUser)
	group.DELETE("/:uid", h.deleteUser)

	group = h.engine.Group("/post")
	group.POST("", h.createPost)
	group.GET("/:uid", h.getPost)
	group.GET("/user/:uid", h.getUserPosts)
	group.GET("/reply/:uid", h.getPostReplies)

	group = h.engine.Group("/edge")
	group.POST("/:source/:edge/:target", h.createEdge)
	group.DELETE("/:source/:edge/:target", h.deleteEdge)

	group = h.engine.Group("/index")
	group.GET("/search/:query", h.search)
	group.PATCH("/hit/:uid", h.hit)

	group = h.engine.Group("/storage")
	group.POST("", h.upload)

	return h.engine.Run(":" + port)
}

func (h *Handler) root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
