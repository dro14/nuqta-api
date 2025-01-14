package handler

import (
	"net/http"
	"strings"

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

	group = authorized.Group("/user")
	group.POST("", h.createUser)
	group.GET("/:by", h.getUser)
	group.PUT("", h.updateUser)
	group.DELETE("/:uid/:predicate", h.deleteUserPredicate)
	group.DELETE("/:uid", h.deleteUser)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getPosts)
	group.GET("/:uid", h.getPost)
	group.GET("/user/:uid", h.getUserPosts)
	group.GET("/reply/:uid", h.getPostReplies)

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
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, failure(e.ErrNoAuthHeader))
		c.Abort()
		return
	}

	idToken := strings.TrimPrefix(header, "Bearer ")

	uid, err := h.auth.VerifyIdToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, failure(err))
		c.Abort()
		return
	}

	c.Set("firebase_uid", uid)
}
