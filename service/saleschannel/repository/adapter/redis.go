package adapter

import (
	"github.com/go-redis/redis"

	"go-poc/service/saleschannel/repository/adapter/channel"
	"go-poc/service/saleschannel/repository/port"
)

type redisRegistry struct {
	db *redis.Client
}

func NewRedis(db *redis.Client) port.CacheRepository {
	return redisRegistry{
		db: db,
	}
}

func (r redisRegistry) Channel() port.ChannelCacheRepository {
	return channel.NewRedisRepository(r.db)
}
