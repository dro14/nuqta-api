package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	client *memcache.Client
}

func New() *Memcached {
	// uri, ok := os.LookupEnv("MEMCACHED_URI")
	// if !ok {
	// 	log.Fatal("memcached uri is not specified")
	// }

	return &Memcached{
		// client: memcache.New(uri),
	}
}
