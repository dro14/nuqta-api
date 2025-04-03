package handler

import (
	"github.com/dro14/nuqta-service/auth"
	"github.com/dro14/nuqta-service/database/dgraph"
	"github.com/dro14/nuqta-service/database/meili"
	"github.com/dro14/nuqta-service/database/memcached"
	"github.com/dro14/nuqta-service/yordamchi"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine    *gin.Engine
	db        *dgraph.Dgraph
	index     *meili.Meili
	cache     *memcached.Memcached
	auth      *auth.Auth
	yordamchi *yordamchi.Yordamchi
}

func New() *Handler {
	return &Handler{
		engine:    gin.Default(),
		db:        dgraph.New(),
		index:     meili.New(),
		cache:     memcached.New(),
		auth:      auth.New(),
		yordamchi: yordamchi.New(),
	}
}
