# Обзор API

### Создание нового кошелька

`POST /api/v1/wallets`  
**Статус ответа**: 201 Created | 500 Internal Server Error  
**Тело ответа**: `uuid` (идентификатор созданного кошелька)

### Получение баланса кошелька

`GET /api/v1/wallets/{wallet_id}`  
**Статус ответа**: 200 OK | 404 Not Found | 500 Internal Server Error  
**Тело ответа**: `int32` (баланс кошелька)

### Пополнение или снятие средств с кошелька

`POST /api/v1/wallets/{wallet_id}`  
**Заголовки запроса**: `"Content-Type": "application/json"`  
**Тело запроса**:

```json
{
	"operation_type": "deposit" | "withdraw",
	"amount": int32
}
```

**Статус ответа**: 200 OK | 400 Bad Request | 404 Not Found | 500 Internal Server Error  
**Тело ответа**: `int32` (новый баланс кошелька)

### Удаление кошелька

`DELETE /api/v1/wallets/{wallet_id}`  
**Статус ответа**: 204 No Content (нет содержимого)

### Проверка состояния сервера

`GET /api/v1/healthz`  
**Статус ответа**: 200 OK | 500 Internal Server Error  
**Тело ответа**: "OK" | "Internal Server Error"

### Получение списка созданных кошельков

`GET /api/v1/wallets`  
**Заголовки запроса**: `Authorization: Basic {base64_encoded_credentials}` (требуется базовая аутентификация)  
**Статус ответа**: 200 OK | 204 No Content | 500 Internal Server Error  
**Тело ответа**:

```json
[
	{
		"id": uuid,
		"balance": int32,
		"created_at": timestamp,
		"updated_at": timestamp
	}
]
```
