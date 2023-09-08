package port

import (
	"github.com/google/uuid"

	"go-poc/service/inventory/model"
)

type SourcingMainRepository interface {
	Create(data *model.Sourcing) error
	Update(data *model.Sourcing) error
	FindByID(id uuid.UUID) (*model.Sourcing, error)
	FindByFilter(filter model.SourcingFilter, lock bool) ([]*model.Sourcing, error)
	FindPage(filter model.SourcingFilter, offset, limit int64) ([]*model.Sourcing, error)
	FindTotalByFilter(filter model.SourcingFilter) (int64, error)
	Delete(filter model.SourcingFilter) error
}

type SourcingCacheRepository interface {
	Set(data *model.Sourcing) error
	Get(id uuid.UUID) (*model.Sourcing, error)
	Delete(id uuid.UUID) error
}
