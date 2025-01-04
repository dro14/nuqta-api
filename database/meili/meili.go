package meili

import (
	"log"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	client meilisearch.ServiceManager
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

	return &Meili{
		client: meilisearch.New(
			uri,
			meilisearch.WithAPIKey(masterKey),
		),
	}
}
