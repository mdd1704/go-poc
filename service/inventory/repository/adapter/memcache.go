package adapter

import (
	"github.com/rainycape/memcache"

	"go-poc/service/inventory/repository/adapter/location"
	"go-poc/service/inventory/repository/adapter/sourcing"
	"go-poc/service/inventory/repository/port"
)

type memcacheRegistry struct {
	db *memcache.Client
}

func NewMemcache(db *memcache.Client) port.CacheRepository {
	return memcacheRegistry{
		db: db,
	}
}

func (r memcacheRegistry) Location() port.LocationCacheRepository {
	return location.NewMemcacheRepository(r.db)
}

func (r memcacheRegistry) Sourcing() port.SourcingCacheRepository {
	return sourcing.NewMemcacheRepository(r.db)
}
