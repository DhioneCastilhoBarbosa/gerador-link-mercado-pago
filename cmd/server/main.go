package main

import (
	"log"
	"os"

	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/database"
	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		_ = godotenv.Load(".env")
	}

	database.Connect()
	r := routes.SetupRouter()
	if err := r.Run(":8090"); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
