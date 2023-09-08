package model

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewLocation(v LocationInput) *Location {
	return &Location{
		ID:        uuid.New(),
		Code:      v.Code,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Location) Update(v LocationInput) {
	m.Code = v.Code
	m.UpdatedAt = time.Now()
}

type LocationInput struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code" binding:"required"`
}

type LocationOutput struct {
	ID      uuid.UUID `json:"id"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
}

type LocationFilter struct {
	IDs   []uuid.UUID `json:"ids"`
	Codes []string    `json:"codes"`
}

type LocationURI struct {
	ID string `uri:"id" binding:"required"`
}
