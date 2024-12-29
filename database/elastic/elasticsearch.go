package elasticsearch

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

type Elasticsearch struct {
	client *elasticsearch.Client
}

func New() *Elasticsearch {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{os.Getenv("ELASTICSEARCH_URI")},
			Username:  os.Getenv("ELASTICSEARCH_USERNAME"),
			Password:  os.Getenv("ELASTICSEARCH_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatal("error creating elasticsearch client ", err)
	}

	return &Elasticsearch{client: client}
}
