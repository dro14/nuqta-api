package handler

import (
	"github.com/dro14/nuqta-service/data"
	"github.com/dro14/nuqta-service/firebase"
	"github.com/dro14/nuqta-service/yordamchi"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine    *gin.Engine
	data      *data.Data
	firebase  *firebase.Firebase
	yordamchi *yordamchi.Yordamchi
}

func New() *Handler {
	return &Handler{
		engine:    gin.Default(),
		data:      data.New(),
		firebase:  firebase.New(),
		yordamchi: yordamchi.New(),
	}
}
