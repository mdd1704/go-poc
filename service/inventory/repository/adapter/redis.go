package adapter

import (
	"github.com/go-redis/redis"

	"go-poc/service/inventory/repository/adapter/location"
	"go-poc/service/inventory/repository/adapter/sourcing"
	"go-poc/service/inventory/repository/port"
)

type redisRegistry struct {
	db *redis.Client
}

func NewRedis(db *redis.Client) port.CacheRepository {
	return redisRegistry{
		db: db,
	}
}

func (r redisRegistry) Location() port.LocationCacheRepository {
	return location.NewRedisRepository(r.db)
}

func (r redisRegistry) Sourcing() port.SourcingCacheRepository {
	return sourcing.NewRedisRepository(r.db)
}
