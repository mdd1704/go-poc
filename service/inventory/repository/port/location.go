package port

import (
	"github.com/google/uuid"

	"go-poc/service/inventory/model"
)

type LocationMainRepository interface {
	Create(data *model.Location) error
	Update(data *model.Location) error
	FindByID(id uuid.UUID) (*model.Location, error)
	FindByFilter(filter model.LocationFilter, lock bool) ([]*model.Location, error)
	FindPage(filter model.LocationFilter, offset, limit int64) ([]*model.Location, error)
	FindTotalByFilter(filter model.LocationFilter) (int64, error)
	Delete(filter model.LocationFilter) error
}

type LocationCacheRepository interface {
	Set(data *model.Location) error
	Get(id uuid.UUID) (*model.Location, error)
	Delete(id uuid.UUID) error
}
