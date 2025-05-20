package models

import (
	"github.com/google/uuid"
)

type Payment struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Titulo       string    `json:"titulo"`
	Valor        float64   `json:"valor"`
	Status       string    `json:"status"`
	InitPoint    string    `json:"init_point"`
	PreferenceID string    `json:"preference_id"`
	USER_ID      string    `json:"user_id"`
}
