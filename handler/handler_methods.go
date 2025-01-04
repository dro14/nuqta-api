package handler

import (
	"errors"
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Run(port string) error {
	h.engine.GET("/", h.Root)

	h.engine.GET("/schema", h.GetSchema)
	h.engine.PUT("/schema", h.SetSchema)
	h.engine.DELETE("/schema", h.DeleteSchema)

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
	schema, err := h.db.ReadSchema(c.Request.Context())
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

func (h *Handler) PostUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	existingUser, err := h.db.ReadUserByFirebaseUid(ctx, user.FirebaseUid)
	if err != nil && !errors.Is(err, e.ErrNotFound) {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusOK, existingUser)
		return
	}

	user, err = h.db.CreateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	var user *models.User
	var err error
	switch {
	case c.Query("uid") != "":
		user, err = h.db.ReadUserByUid(ctx, c.Query("uid"))
	case c.Query("firebase_uid") != "":
		user, err = h.db.ReadUserByFirebaseUid(ctx, c.Query("firebase_uid"))
	default:
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) PutUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	user, err = h.db.UpdateUser(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) PatchUser(c *gin.Context) {}

func (h *Handler) DeleteUser(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	err := h.db.DeleteUser(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SearchUser(c *gin.Context) {}
