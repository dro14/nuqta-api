package recommender

import (
	"log"

	"github.com/dro14/nuqta-service/database/dgraph"
	"github.com/dro14/nuqta-service/database/memcached"
	"github.com/dro14/nuqta-service/models"
)

type Recommender struct {
	db    *dgraph.Dgraph
	cache *memcached.Memcached
	recs  []*models.Post
}

func New() *Recommender {
	cache := memcached.New()
	recs, err := cache.GetRecs()
	if err != nil {
		log.Println(err)
	}

	return &Recommender{
		db:    dgraph.New(),
		cache: cache,
		recs:  recs,
	}
}
