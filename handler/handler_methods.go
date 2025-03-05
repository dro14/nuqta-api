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
	group.DELETE("/:predicate", h.deletePredicate)

	authorized := h.engine.Group("")
	authorized.Use(h.authMiddleware)

	group = authorized.Group("/profile")
	group.POST("", h.createProfile)
	group.GET("", h.getProfile)
	group.PUT("", h.updateProfile)
	group.PATCH("", h.updateProfileAttribute)
	group.DELETE("", h.deleteProfileAttribute)
	group.OPTIONS("", h.isAvailable)

	group = authorized.Group("/user")
	group.GET("", h.getUser)
	group.GET("/followers", h.getUserFollowers)
	group.GET("/following", h.getUserFollowing)
	group.GET("/search", h.searchUser)
	group.PATCH("", h.hitUser)

	group = authorized.Group("/post")
	group.POST("", h.createPost)
	group.GET("", h.getPost)
	group.DELETE("", h.deletePost)

	group = authorized.Group("/edge")
	group.POST("", h.createEdge)
	group.DELETE("", h.deleteEdge)

	return h.engine.Run(":" + port)
}

func (h *Handler) UpdateRecs() {
	go h.rec.UpdateRecs()
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
