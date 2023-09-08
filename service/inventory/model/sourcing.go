package model

import (
	"time"

	"github.com/google/uuid"
)

type Sourcing struct {
	ID          uuid.UUID `json:"id" db:"id"`
	SKU         string    `json:"sku" db:"sku"`
	QtyTotal    int       `json:"qty_total" db:"qty_total"`
	QtyReserved int       `json:"qty_reserved" db:"qty_reserved"`
	QtySaleable int       `json:"qty_saleable" db:"qty_saleable"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewSourcing(v SourcingInput) *Sourcing {
	return &Sourcing{
		ID:          uuid.New(),
		SKU:         v.SKU,
		QtyTotal:    v.QtyTotal,
		QtyReserved: v.QtyReserved,
		QtySaleable: v.QtySaleable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (m *Sourcing) Update(v SourcingInput) {
	m.QtyTotal = v.QtyTotal
	m.QtyReserved = v.QtyReserved
	m.QtySaleable = v.QtySaleable
	m.UpdatedAt = time.Now()
}

type SourcingInput struct {
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku" binding:"required"`
	QtyTotal    int       `json:"qty_total"`
	QtyReserved int       `json:"qty_reserved"`
	QtySaleable int       `json:"qty_saleable"`
}

type SourcingOutput struct {
	ID      uuid.UUID `json:"id"`
	SKU     string    `json:"sku"`
	Message string    `json:"message"`
}

type SourcingFilter struct {
	IDs  []uuid.UUID `json:"ids"`
	SKUs []string    `json:"skus"`
}

type SourcingURI struct {
	ID string `uri:"id" binding:"required"`
}
