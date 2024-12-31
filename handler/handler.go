package handler

import (
	"github.com/dro14/nuqta-service/database/dgraph"
	"github.com/dro14/nuqta-service/database/elastic"
	"github.com/dro14/nuqta-service/database/memcached"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine *gin.Engine
	db     *dgraph.Dgraph
	search *elastic.Elastic
	cache  *memcached.Memcached
}

func New() *Handler {
	return &Handler{
		engine: gin.Default(),
		db:     dgraph.New(),
		search: elastic.New(),
		cache:  memcached.New(),
	}
}
