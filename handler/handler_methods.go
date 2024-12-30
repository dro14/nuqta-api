package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Print("can't bind json: ", err)
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := context.WithValue(c.Request.Context(), "id", user.ID)

	err = h.mongo.CreateUser(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		c.JSON(http.StatusCreated, nil)
		return
	} else if err != nil {
		log.Print("can't create user in mongo")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.elastic.CreateUser(ctx, user)
	if err != nil {
		log.Print("can't create user in elastic")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, failure(errNoID))
		return
	}
	ctx := context.WithValue(c.Request.Context(), "id", id)

	user, err := h.mongo.ReadUser(ctx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user")
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err != nil {
		log.Print("can't read user")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *Handler) PutUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Print("can't bind json: ", err)
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := context.WithValue(c.Request.Context(), "id", user.ID)

	err = h.mongo.UpdateUser(ctx, user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user")
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err != nil {
		log.Print("can't update user")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) PatchUser(c *gin.Context) {

}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, failure(errNoID))
		return
	}
	ctx := context.WithValue(c.Request.Context(), "id", id)

	err := h.mongo.DeleteUser(ctx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user")
		c.JSON(http.StatusNotFound, nil)
		return
	} else if err != nil {
		log.Print("can't delete user")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) SearchUser(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, failure(errNoQuery))
		return
	}
	ctx := context.WithValue(c.Request.Context(), "query", query)

	ids, err := h.elastic.SearchUser(ctx, query)
	if err != nil {
		log.Print("can't search user")
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"ids": ids})
}
