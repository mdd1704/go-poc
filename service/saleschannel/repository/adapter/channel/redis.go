package channel

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	"go-poc/service/saleschannel/model"
	"go-poc/service/saleschannel/repository/port"
)

type redisRepository struct {
	db *redis.Client
}

func NewRedisRepository(db *redis.Client) port.ChannelCacheRepository {
	return &redisRepository{
		db: db,
	}
}

func (repo *redisRepository) Set(data *model.Channel) error {
	value, err := json.Marshal(*data)
	if err != nil {
		return err
	}

	result := repo.db.Set(data.ID.String(), string(value), time.Duration(time.Hour*24*30))
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (repo *redisRepository) Get(id uuid.UUID) (data *model.Channel, err error) {
	result, err := repo.db.Get(id.String()).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("not found")
	}

	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *redisRepository) Delete(id uuid.UUID) error {
	result := repo.db.Del(id.String())
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
