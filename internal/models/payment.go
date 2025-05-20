package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Titulo       string  `json:"titulo"`
	Valor        float64 `json:"valor"`
	Status       string  `json:"status"` // pending, paid, failed
	InitPoint    string  `json:"init_point"`
	PreferenceID string  `json:"preference_id"`
}
