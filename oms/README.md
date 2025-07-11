# OMS (Order Management System)

Welcome to OMS, a Go-based microservice for robust order management. OMS is designed for bulk order processing, CSV uploads, and seamless integration with external services like S3, SQS, Kafka, MongoDB, and your IMS (Inventory Management System).

---

## ✨ Features

- **Bulk Order Processing:** Upload and process CSV files containing order data.
- **CSV Validation:** Ensures uploaded files have all required columns and valid data.
- **Multi-Service Integration:** Works with S3 (file storage), SQS (queueing), Kafka (event streaming), and MongoDB (order storage).
- **Asynchronous Processing:** Orders are processed in the background for scalability.
- **Inventory Validation:** Real-time checks with IMS before order creation.
- **Internationalization:** Multi-language support via i18n.
- **Structured Logging & Monitoring:** Built-in logging and health checks.
- **Docker Support:** Easy local development with Docker Compose.

---

## 🗂️ Project Structure

```
oms/
├── configs/         # YAML configuration (if used)
├── constants/       # Application constants
├── controllers/     # HTTP handlers (Gin)
├── database/        # MongoDB connection and repositories
├── middleware/      # Auth and logging middleware
├── models/          # Data models
├── routes/          # Route registration
├── utils/           # S3, SQS, Kafka, CSV, IMS client, etc.
├── webhook/         # Webhook event handling
├── localization/    # i18n files
├── main.go          # Service entry point
├── go.mod, go.sum   # Go dependencies
```

---

## 🚦 Prerequisites

- Go 1.20 or newer
- MongoDB (for order storage)
- Apache Kafka (for event streaming)
- AWS S3 & SQS (or LocalStack for local dev)
- IMS Service (for inventory validation)
- Docker & Docker Compose (recommended for local development)

---

## ⚙️ Configuration

Set environment variables (or use a `.env` file):

```env
# MongoDB
MONGODB_URI=mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin
MONGODB_DB_NAME=mydb

# AWS/LocalStack
AWS_REGION=us-east-1
AWS_S3_ENDPOINT=http://localhost:4566
AWS_SQS_ENDPOINT=http://localhost:4566
S3_BUCKET_NAME=order-csv-bucket
CREATE_BULK_ORDER_QUEUE_NAME=CreateBulkOrder

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_ORDER_TOPIC=order.created

# IMS Service
IMS_BASE_URL=http://localhost:8084
```

---

## 🚀 Running the Service

### With Docker Compose (Recommended)

```sh
docker-compose up -d
```
- The service will be available at: [http://localhost:8086](http://localhost:8086)

### For Local Development

```sh
go mod download
go run main.go
```

---

## 📚 API Endpoints

### Orders

| Method | Endpoint                        | Description                        |
|--------|---------------------------------|------------------------------------|
| POST   | /api/v1/orders/upload           | Upload CSV for bulk order creation |
| GET    | /api/v1/orders/                 | List orders (with filters)         |
| GET    | /api/v1/orders/:orderID         | Get order by ID                    |
| PUT    | /api/v1/orders/:orderID/status  | Update order status                |

### Webhooks

| Method | Endpoint                  | Description                |
|--------|---------------------------|----------------------------|
| GET    | /api/v1/webhook/events    | List webhook events        |

---

## 🧩 Data Models

### Order

```go
type Order struct {
    ID        string      `json:"id"`
    SKU       string      `json:"sku"`
    Location  string      `json:"location"`
    TenantID  string      `json:"tenant_id"`
    SellerID  string      `json:"seller_id"`
    Status    string      `json:"status"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
}
```

### WebhookEvent

```go
type WebhookEvent struct {
    ID        string    `json:"id"`
    EventType string    `json:"event_type"`
    Payload   string    `json:"payload"`
    CreatedAt time.Time `json:"created_at"`
}
```

---

## 🛠️ Integration with IMS

OMS validates inventory and SKU/hub existence by calling your IMS service before creating orders.  
IMS endpoints used:
- `/sku/`
- `/hub/`
- `/inventory/`

---

## 🧪 Integration Testing

- Integration tests are in `controllers/*_test.go`, `utils/*_test.go`, and `webhook/*_test.go`.
- Run all tests with:
  ```sh
  go test ./...
  ```
- For real integration, ensure MongoDB, Kafka, S3/SQS (or LocalStack), and IMS are running.

---

## 💡 Notes

- All endpoints expect and return JSON.
- Use query parameters for filtering orders (e.g., `tenant_id`, `seller_id`, `status`, `page`, `limit`).
- Redis is not used in OMS, but S3, SQS, and Kafka are required for full functionality.
- The `/webhook/events` HTML page is available at [http://localhost:8086/webhook/events](http://localhost:8086/webhook/events) for viewing webhook events.

---

## 🤝 Contributing

- Follow Go best practices.
- Add tests for new features.
- Update documentation for any API changes.

---

If you have questions or want to see usage examples, feel free to ask or open an issue!