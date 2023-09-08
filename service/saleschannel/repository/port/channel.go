package port

import (
	"github.com/google/uuid"

	"go-poc/service/saleschannel/model"
)

type ChannelMainRepository interface {
	Create(data *model.Channel) error
	Update(data *model.Channel) error
	FindByID(id uuid.UUID) (*model.Channel, error)
	FindByFilter(filter model.ChannelFilter, lock bool) ([]*model.Channel, error)
	FindPage(filter model.ChannelFilter, offset, limit int64) ([]*model.Channel, error)
	FindTotalByFilter(filter model.ChannelFilter) (int64, error)
	Delete(filter model.ChannelFilter) error
}

type ChannelCacheRepository interface {
	Set(data *model.Channel) error
	Get(id uuid.UUID) (*model.Channel, error)
	Delete(id uuid.UUID) error
}
