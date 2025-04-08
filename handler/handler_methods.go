package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.root)

	group := h.engine.Group("/client")
	group.GET("", h.getClientInfo)

	group = h.engine.Group("/schema")
	group.GET("", h.getSchema)
	group.PUT("", h.updateSchema)
	group.DELETE("", h.deleteSchema)
	group.DELETE("/:predicate", h.deletePredicate)

	authorized := h.engine.Group("")
	authorized.Use(h.authMiddleware)

	group = authorized.Group("/profile")
	group.POST("", h.createProfile)
	group.GET("", h.getProfile)
	group.PUT("", h.updateProfile)
	group.OPTIONS("", h.isAvailable)

	group = authorized.Group("/user")
	group.GET("", h.getUserList)
	group.GET("/username", h.getUserByUsername)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getPostList)
	group.DELETE("", h.deletePost)

	group = authorized.Group("/edge")
	group.POST("", h.createEdge)
	group.DELETE("", h.deleteEdge)

	group = authorized.Group("/yordamchi")
	group.POST("", h.createResponse)

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

	firebaseUid, err := h.auth.VerifyIdToken(c.Request.Context(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, failure(err))
		c.Abort()
		return
	}

	c.Set("firebase_uid", firebaseUid)
}

func (h *Handler) getClientInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ip_address": c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
	})
}
