package handler

import (
	"github.com/dro14/nuqta-service/database/mongo"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
	mongo  *mongo.Mongo
}

func New() *Handler {
	return &Handler{
		engine: gin.Default(),
		mongo:  mongo.New(),
	}
}
