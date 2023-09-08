package model

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewChannel(v ChannelInput) *Channel {
	return &Channel{
		ID:        uuid.New(),
		Code:      v.Code,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Channel) Update(v ChannelInput) {
	m.Code = v.Code
	m.UpdatedAt = time.Now()
}

type ChannelInput struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code" binding:"required"`
}

type ChannelOutput struct {
	ID      uuid.UUID `json:"id"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
}

type ChannelFilter struct {
	IDs   []uuid.UUID `json:"ids"`
	Codes []string    `json:"codes"`
}

type ChannelURI struct {
	ID string `uri:"id" binding:"required"`
}
