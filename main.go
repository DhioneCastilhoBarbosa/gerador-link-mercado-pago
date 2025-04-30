package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type PagamentoRequest struct {
	Titulo string  `json:"titulo"`
	Valor  float64 `json:"valor"`
}

func main() {
	if os.Getenv("GIN_MODE") != "release" {
		_ = godotenv.Load(".env") // Carrega .env somente em dev
	}

	accessToken := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")
	if accessToken == "" {
		panic("MERCADO_PAGO_ACCESS_TOKEN não encontrado")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // você pode restringir para ["https://www.seusite.com"]
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.POST("/criar-pagamento", func(c *gin.Context) {
		var req PagamentoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
			return
		}

		preference := map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"title":      req.Titulo,
					"quantity":   1,
					"unit_price": req.Valor,
				},
			},
			"back_urls": map[string]string{
				"success": "https://www.eletrihub.com/",
				"failure": "https://sua-url.com/erro",
				"pending": "https://sua-url.com/pendente",
			},
			"auto_return": "approved",
		}

		body, _ := json.Marshal(preference)

		reqMercadoPago, _ := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewBuffer(body))
		reqMercadoPago.Header.Set("Content-Type", "application/json")
		reqMercadoPago.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		resp, err := client.Do(reqMercadoPago)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pagamento"})
			return
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)

		var respJson map[string]interface{}
		json.Unmarshal(respBody, &respJson)

		if initPoint, ok := respJson["init_point"].(string); ok {
			c.JSON(http.StatusOK, gin.H{"url": initPoint})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter link de pagamento"})
		}
	})

	r.Run(":8088")
}
