package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Titulo       string    `json:"titulo"`
	Valor        float64   `json:"valor"`
	Status       string    `json:"status"`
	InitPoint    string    `json:"init_point"`
	PreferenceID string    `json:"preference_id"`
}

// BeforeCreate gera UUID automaticamente
func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
