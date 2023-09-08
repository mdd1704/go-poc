package adapter

import (
	"github.com/rainycape/memcache"

	"go-poc/service/saleschannel/repository/adapter/channel"
	"go-poc/service/saleschannel/repository/port"
)

type memcacheRegistry struct {
	db *memcache.Client
}

func NewMemcache(db *memcache.Client) port.CacheRepository {
	return memcacheRegistry{
		db: db,
	}
}

func (r memcacheRegistry) Channel() port.ChannelCacheRepository {
	return channel.NewMemcacheRepository(r.db)
}
