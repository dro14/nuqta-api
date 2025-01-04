package meili

import (
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	index meilisearch.IndexManager
}

func New() *Meili {
	uri, ok := os.LookupEnv("MEILI_URI")
	if !ok {
		log.Fatal("meili uri is not specified")
	}

	masterKey, ok := os.LookupEnv("MEILI_MASTER_KEY")
	if !ok {
		log.Fatal("meili master key is not specified")
	}

	client := meilisearch.New(uri, meilisearch.WithAPIKey(masterKey))
	index := client.Index("users")
	index.UpdateSortableAttributes(&[]string{"hits"})
	return &Meili{
		index: index,
	}
}
