package handler

import (
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

	err = h.mongo.CreateUser(c.Request.Context(), user)
	if mongo.IsDuplicateKeyError(err) {
		err = h.mongo.UpdateUser(c.Request.Context(), user)
		if err != nil {
			log.Print("can't update user: ", err)
			c.JSON(http.StatusInternalServerError, failure(err))
		} else {
			c.JSON(http.StatusNoContent, nil)
		}
	} else if err != nil {
		log.Print("can't create user: ", err)
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		c.JSON(http.StatusCreated, nil)
	}
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, failure(errNoID))
		return
	}

	user, err := h.mongo.ReadUser(c.Request.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user: ", err)
		c.JSON(http.StatusNotFound, failure(err))
	} else if err != nil {
		log.Print("can't read user: ", err)
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func (h *Handler) PutUser(c *gin.Context) {
	user := &models.User{}
	err := c.ShouldBindJSON(user)
	if err != nil {
		log.Print("can't bind json: ", err)
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	err = h.mongo.UpdateUser(c.Request.Context(), user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user: ", err)
		c.JSON(http.StatusNotFound, failure(err))
	} else if err != nil {
		log.Print("can't update user: ", err)
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}

func (h *Handler) PatchUser(c *gin.Context) {

}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, failure(errNoID))
		return
	}

	err := h.mongo.DeleteUser(c.Request.Context(), id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Print("can't find user: ", err)
		c.JSON(http.StatusNotFound, failure(err))
	} else if err != nil {
		log.Print("can't delete user: ", err)
		c.JSON(http.StatusInternalServerError, failure(err))
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
