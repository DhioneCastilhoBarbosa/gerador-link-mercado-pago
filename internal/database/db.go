package database

import (
	"log"
	"os"

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

	// Executa migrações com versionamento
	if err := RunMigrations(DB); err != nil {
		log.Fatal("Erro ao rodar migrações:", err)
	}
}
