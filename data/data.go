package data

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	_ "github.com/lib/pq"
)

type Data struct {
	db    *sql.DB
	graph *dgo.Dgraph
	cache *memcache.Client
}

func New() *Data {
	uri, ok := os.LookupEnv("POSTGRES_URI")
	if !ok {
		log.Fatal("postgres uri is not specified")
	}

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("can't connect to postgres: ", err)
	}

	uri, ok = os.LookupEnv("DGRAPH_URI")
	if !ok {
		log.Fatal("dgraph uri is not specified")
	}

	graph, err := dgo.Open(uri)
	if err != nil {
		log.Fatal("can't connect to dgraph: ", err)
	}

	uri, ok = os.LookupEnv("MEMCACHED_URI")
	if !ok {
		log.Fatal("memcached uri is not specified")
	}

	cache := memcache.New(uri)

	return &Data{
		db:    db,
		graph: graph,
		cache: cache,
	}
}

func (d *Data) DeleteType(ctx context.Context, name string) error {
	op := &api.Operation{
		DropOp:    api.Operation_TYPE,
		DropValue: name,
	}
	return d.graph.Alter(ctx, op)
}
