package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/database"
	"github.com/DhioneCastilhoBarbosa/mercado-pago-api-link/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// Gera UUID para servir como ID e external_reference
	id := uuid.New()

	preference := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"title":      req.Titulo,
				"quantity":   1,
				"unit_price": req.Valor,
			},
		},
		"external_reference": id.String(),
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

	initPoint, _ := respJson["init_point"].(string)
	preferenceID, _ := respJson["id"].(string)

	p := models.Payment{
		ID:           id,
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
	topic := c.Query("topic")
	if topic == "" {
		topic = c.Query("type")
	}

	id := c.Query("id")
	if id == "" {
		id = c.Query("data.id")
	}

	if topic != "payment" || id == "" {
		log.Printf("🔔 Webhook ignorado: topic=%s, id=%s\n", topic, id)
		c.JSON(http.StatusOK, gin.H{"message": "Notificação ignorada"})
		return
	}

	log.Printf("✅ Webhook recebido: topic=%s, id=%s\n", topic, id)

	token := os.Getenv("MERCADO_PAGO_ACCESS_TOKEN")
	url := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s?access_token=%s", id, token)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("❌ Erro ao consultar pagamento:", err)
		c.JSON(http.StatusOK, gin.H{"message": "Erro na consulta"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("📦 Resposta completa do pagamento:\n%s\n", string(body))

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Println("❌ Erro ao decodificar JSON da resposta:", err)
		c.JSON(http.StatusOK, gin.H{"message": "Erro na decodificação"})
		return
	}

	externalRef, _ := raw["external_reference"].(string)
	status, _ := raw["status"].(string)

	if externalRef == "" || status == "" {
		log.Println("⚠️ external_reference ou status ausente na resposta")
		c.JSON(http.StatusOK, gin.H{"message": "Dados incompletos"})
		return
	}

	var pagamento models.Payment
	if err := database.DB.Where("id = ?", externalRef).First(&pagamento).Error; err != nil {
		log.Printf("❌ Pagamento com ID %s não encontrado", externalRef)
		c.JSON(http.StatusOK, gin.H{"message": "Pagamento não encontrado"})
		return
	}

	pagamento.Status = status
	database.DB.Save(&pagamento)

	log.Printf("✅ Pagamento atualizado: %s → %s", pagamento.ID.String(), status)
	c.Status(http.StatusOK)
}
