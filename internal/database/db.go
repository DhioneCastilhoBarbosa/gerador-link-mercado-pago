package database

import (
	"log"
	"os"

	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar no banco de dados:", err)
	}

	if err := DB.AutoMigrate(&models.Payment{}); err != nil {
		log.Fatal("Erro ao realizar migração:", err)
	}
}
