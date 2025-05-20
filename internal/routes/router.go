package routes

import (
	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.POST("/criar-pagamento", handlers.CriarPagamento)
	r.POST("/webhook-mercado-pago", handlers.WebhookMercadoPago)

	return r
}
