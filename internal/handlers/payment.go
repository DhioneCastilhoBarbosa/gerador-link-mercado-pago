package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/database"
	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/models"
	"github.com/gin-gonic/gin"
)

type PagamentoRequest struct {
	Titulo string  `json:"titulo"`
	Valor  float64 `json:"valor"`
}

func CriarPagamento(c *gin.Context) {
	var req PagamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
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
		"notification_url": "https://api.eletrihub.com/webhook-mercado-pago",
		"auto_return":      "approved",
	}

	body, _ := json.Marshal(preference)

	reqMP, _ := http.NewRequest("POST", "https://api.mercadopago.com/checkout/preferences", bytes.NewBuffer(body))
	reqMP.Header.Set("Content-Type", "application/json")
	reqMP.Header.Set("Authorization", "Bearer "+os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(reqMP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pagamento"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var respJson map[string]interface{}
	json.Unmarshal(respBody, &respJson)

	initPoint := respJson["init_point"].(string)
	preferenceID := respJson["id"].(string)

	p := models.Payment{
		Titulo:       req.Titulo,
		Valor:        req.Valor,
		Status:       "pending",
		InitPoint:    initPoint,
		PreferenceID: preferenceID,
	}
	database.DB.Create(&p)

	c.JSON(http.StatusOK, gin.H{"url": initPoint})
}

func WebhookMercadoPago(c *gin.Context) {
	idPagamento := c.Query("id")
	if idPagamento == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do pagamento não informado"})
		return
	}

	token := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s?access_token=%s", idPagamento, token)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar pagamento"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var pagamentoMP struct {
		Status       string `json:"status"`
		PreferenceID string `json:"order"`
	}

	_ = json.Unmarshal(body, &pagamentoMP)

	if pagamentoMP.PreferenceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pagamento não encontrado"})
		return
	}

	var pagamento models.Payment
	if err := database.DB.Where("preference_id = ?", pagamentoMP.PreferenceID).First(&pagamento).Error; err == nil {
		pagamento.Status = pagamentoMP.Status
		database.DB.Save(&pagamento)
	}

	c.Status(http.StatusOK)
}
