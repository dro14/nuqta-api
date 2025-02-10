package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.root)

	group := h.engine.Group("/schema")
	group.GET("", h.getSchema)
	group.PUT("", h.updateSchema)
	group.DELETE("", h.deleteSchema)

	authorized := h.engine.Group("")
	authorized.Use(h.authMiddleware)

	group = authorized.Group("/profile")
	group.POST("", h.createProfile)
	group.GET("", h.getProfile)
	group.PUT("", h.updateProfile)
	group.DELETE("/:uid/:attribute", h.deleteProfileAttribute)

	group = authorized.Group("/user")
	group.GET("", h.getUser)
	group.GET("/:username", h.isUsernameAvailable)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getAllPosts)
	group.GET("/:uid", h.getPost)
	group.GET("/following/:before", h.getFollowingPosts)
	group.GET("/user/:uid", h.getUserPosts)
	group.GET("/reply/:uid", h.getPostReplies)
	group.DELETE("/:uid", h.deletePost)

	group = authorized.Group("/edge")
	group.POST("/:source/:edge/:target", h.createEdge)
	group.DELETE("/:source/:edge/:target", h.deleteEdge)

	group = authorized.Group("/index")
	group.GET("/search", h.search)
	group.PATCH("/hit/:uid", h.hit)

	group = authorized.Group("/storage")
	group.POST("", h.upload)

	return h.engine.Run(":" + port)
}

func (h *Handler) root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (h *Handler) authMiddleware(c *gin.Context) {
	idToken := c.GetHeader("my-id-token")
	if idToken == "" {
		c.JSON(http.StatusUnauthorized, failure(e.ErrNoAuthHeader))
		c.Abort()
		return
	}

	uid, err := h.auth.VerifyIdToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, failure(err))
		c.Abort()
		return
	}

	c.Set("firebase_uid", uid)
}
