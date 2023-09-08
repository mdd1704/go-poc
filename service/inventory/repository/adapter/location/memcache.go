package location

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rainycape/memcache"

	"go-poc/service/inventory/model"
	"go-poc/service/inventory/repository/port"
)

type memcacheRepository struct {
	db *memcache.Client
}

func NewMemcacheRepository(db *memcache.Client) port.LocationCacheRepository {
	return &memcacheRepository{
		db: db,
	}
}

func (repo *memcacheRepository) Set(data *model.Location) error {
	dataMarshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = repo.db.Set(&memcache.Item{Key: data.ID.String(), Value: dataMarshal})
	if err != nil {
		return err
	}

	return nil
}

func (repo *memcacheRepository) Get(id uuid.UUID) (data *model.Location, err error) {
	result, err := repo.db.Get(id.String())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result.Value), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *memcacheRepository) Delete(id uuid.UUID) error {
	err := repo.db.Delete(id.String())
	if err != nil {
		return err
	}

	return nil
}
