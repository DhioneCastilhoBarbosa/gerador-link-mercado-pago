# API de Pagamento com Mercado Pago

Esta API foi desenvolvida em Go usando o framework Gin e permite gerar links de pagamento através do Mercado Pago. O link gerado pode ser usado no frontend para redirecionar o usuário à finalização do pagamento.

---

## 🚀 Funcionalidades

- Recebe o valor e o título do produto via requisição POST
- Cria uma preferência de pagamento no Mercado Pago
- Retorna a URL (`init_point`) para redirecionamento do usuário
- Define URLs de retorno para sucesso, erro ou pendência
- Suporte a CORS para integração com frontends

---

## 🧪 Teste com Postman

**Endpoint:**

```
POST http://localhost:8088/criar-pagamento
```

**Body (JSON):**

```json
{
  "titulo": "Produto Exemplo",
  "valor": 99.90
}
```

**Resposta esperada:**

```json
{
  "url": "https://www.mercadopago.com.br/checkout/v1/redirect?preference_id=..."
}
```

---

## ⚙️ Como executar localmente

### 1. Clone o projeto

```bash
git clone https://github.com/seu-usuario/sua-api-pagamento.git
cd sua-api-pagamento
```

### 2. Crie um arquivo `.env`

```env
MERCADO_PAGO_ACCESS_TOKEN=seu_token_do_mercado_pago
```

> **Importante**: o arquivo `.env` é usado apenas em desenvolvimento.

### 3. Instale as dependências

```bash
go mod tidy
```

### 4. Rode o servidor

```bash
go run main.go
```

A API estará disponível em `http://localhost:8088`.

---

## 🚀 Como rodar em produção

Em produção, as variáveis de ambiente devem ser definidas diretamente no ambiente e o `.env` **não será carregado**.

Exemplo de execução:

```bash
GIN_MODE=release MERCADO_PAGO_ACCESS_TOKEN=seu_token go run main.go
```

Você também pode usar um gerenciador de processos como `systemd`, `pm2`, `docker` etc.

---

## 📦 Dependências

- [Gin](https://github.com/gin-gonic/gin)
- [Mercado Pago API](https://www.mercadopago.com.br/developers/pt)
- [godotenv](https://github.com/joho/godotenv)
- [gin-contrib/cors](https://github.com/gin-contrib/cors)

---

## 🛡️ Boas práticas

- Nunca suba seu `.env` para o repositório (adicione ao `.gitignore`)
- Em produção, sempre use `GIN_MODE=release`
- Gere tokens seguros no painel do Mercado Pago

---

## 📄 Licença

Este projeto está licenciado sob a Licença MIT.
