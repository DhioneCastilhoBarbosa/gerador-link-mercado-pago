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
		c.JSON(http.StatusOK, gin.H{"message": "Erro na consulta, mas webhook aceito para não gerar retry em loop"})
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

	// Extrai dados com segurança
	preferenceID, _ := raw["preference_id"].(string)
	status, _ := raw["status"].(string)

	if preferenceID == "" || status == "" {
		log.Println("⚠️ Dados de pagamento incompletos na resposta do Mercado Pago")
		c.JSON(http.StatusOK, gin.H{"message": "Pagamento sem dados suficientes"})
		return
	}

	// Atualiza no banco
	var pagamento models.Payment
	if err := database.DB.Where("preference_id = ?", preferenceID).First(&pagamento).Error; err != nil {
		log.Printf("❌ Pagamento não encontrado no banco: preference_id=%s", preferenceID)
		c.JSON(http.StatusOK, gin.H{"message": "Pagamento não encontrado, mas notificação recebida"})
		return
	}

	pagamento.Status = status
	database.DB.Save(&pagamento)

	log.Printf("✅ Pagamento atualizado: %s → %s", preferenceID, status)
	c.Status(http.StatusOK)
}
