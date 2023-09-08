package external

import (
	"os"

	"github.com/rainycape/memcache"
)

func NewMemcache() (*memcache.Client, error) {
	dbHost := os.Getenv("MEMCACHE_HOST")
	dbPort := os.Getenv("MEMCACHE_PORT")
	return memcache.New(dbHost + ":" + dbPort)
}
