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
	group.PATCH("/follow/:follower_uid/:followee_uid", h.followUser)
	group.PATCH("/unfollow/:follower_uid/:followee_uid", h.unfollowUser)
	group.DELETE("/:uid", h.deleteUser)

	// group = h.engine.Group("/post")
	// TODO: add post methods

	group = h.engine.Group("/index")
	group.GET("/search", h.search)
	group.PATCH("/increment/:uid", h.increment)

	group = h.engine.Group("/storage")
	group.POST("", h.upload)

	return h.engine.Run(":" + port)
}

func (h *Handler) root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
