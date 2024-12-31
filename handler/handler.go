package handler

import (
	"github.com/dro14/nuqta-service/database/dgraph"
	"github.com/dro14/nuqta-service/database/elastic"
	"github.com/dro14/nuqta-service/database/memcached"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	engine    *gin.Engine
	dgraph    *dgraph.Dgraph
	elastic   *elastic.Elastic
	memcached *memcached.Memcached
}

func New() *Handler {
	return &Handler{
		engine:    gin.Default(),
		elastic:   elastic.New(),
		dgraph:    dgraph.New(),
		memcached: memcached.New(),
	}
}
