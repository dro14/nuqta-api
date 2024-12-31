package dgraph

import (
	"log"
	"os"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Dgraph struct {
	client *dgo.Dgraph
}

func New() *Dgraph {
	uri, ok := os.LookupEnv("DGRAPH_URI")
	if !ok {
		log.Fatal("dgraph uri is not specified")
	}

	conn, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("can't connect to dgraph")
	}

	return &Dgraph{
		client: dgo.NewDgraphClient(api.NewDgraphClient(conn)),
	}
}
