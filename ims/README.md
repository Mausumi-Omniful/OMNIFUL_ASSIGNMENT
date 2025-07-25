# IMS (Inventory Management System)

A Go-based microservice for inventory management, supporting multi-tenant SKU, hub, and inventory operations. Built with Gin, GORM, PostgreSQL, Redis, and the `omniful/go_commons` library.

---

## Features

- **SKU Management:** CRUD for product SKUs, with tenant and seller filtering.
- **Hub Management:** CRUD for warehouse locations, with tenant and seller filtering.
- **Inventory Management:** CRUD, upsert, and atomic reduction of inventory levels.
- **Multi-tenant Support:** All entities support `tenant_id` and `seller_id`.
- **Caching:** Redis caching for fast SKU, hub, and inventory lookups.
- **Database Migrations:** Automatic schema migrations on startup.
- **RESTful API:** Clean, versioned endpoints.
- **Docker Support:** Easy local development with Docker Compose.
- **YAML-based Configuration:** Centralized config in `ims/configs/config.yaml`.

---

## Project Structure

```
ims/
├── configs/         # YAML configuration
├── constants/       # Application constants
├── controllers/     # HTTP handlers (Gin)
├── db/              # Database connection and migrations
├── migrations/      # SQL migration files
├── models/          # Data models
├── redisclient/     # Redis client utilities
├── routes/          # Route registration
├── main.go          # Service entry point
├── go.mod, go.sum   # Go dependencies
```

---

## Prerequisites

- Go 1.20+
- PostgreSQL 13+
- Redis 7.0+
- Docker & Docker Compose (recommended for local dev)

---

## Configuration

Edit `ims/configs/config.yaml`:

```yaml
DB_HOST: localhost
DB_PORT: 5434
DB_USER: postgres
DB_PASSWORD: postgres
DB_NAME: postgres
REDIS_HOST: localhost
REDIS_PORT: 6380
```

---

## Running the Service

### With Docker Compose (Recommended)

```sh
docker-compose up -d
```
- Service: http://localhost:8084

### Local Development

```sh
go mod download
go run main.go
```

---

## Database Migrations

Migrations run automatically on startup.  
To run manually:

```sh
go run db/migration.go
# To rollback:
go run db/migration.go rollback
```

---

## API Endpoints

### Hub

| Method | Endpoint      | Description         |
|--------|--------------|---------------------|
| POST   | /hub/        | Create hub          |
| GET    | /hub/        | List hubs           |
| PUT    | /hub/:id     | Update hub          |
| DELETE | /hub/:id     | Delete hub          |

### SKU

| Method | Endpoint      | Description         |
|--------|--------------|---------------------|
| POST   | /sku/        | Create SKU          |
| GET    | /sku/        | List SKUs           |
| PUT    | /sku/:id     | Update SKU          |
| DELETE | /sku/:id     | Delete SKU          |

### Inventory

| Method | Endpoint           | Description                |
|--------|-------------------|----------------------------|
| POST   | /inventory/       | Create inventory item      |
| GET    | /inventory/       | List inventory             |
| PUT    | /inventory/:id    | Update inventory           |
| DELETE | /inventory/:id    | Delete inventory           |
| POST   | /inventory/upsert | Upsert by SKU/location     |
| POST   | /inventory/reduce | Atomically reduce quantity |

---

## Data Models

### Hub

```go
type Hub struct {
    ID       uint   `json:"id"`
    Name     string `json:"name"`
    Location string `json:"location"`
    TenantID string `json:"tenant_id"`
    SellerID string `json:"seller_id"`
}
```

### SKU

```go
type SKU struct {
    ID          uint   `json:"id"`
    Code        string `json:"sku_code"`
    Name        string `json:"name"`
    Description string `json:"description"`
    TenantID    string `json:"tenant_id"`
    SellerID    string `json:"seller_id"`
}
```

### Inventory

```go
type Inventory struct {
    ID        uint   `json:"id"`
    ProductID string `json:"product_id"`
    SKU       string `json:"sku"`
    Location  string `json:"location"`
    TenantID  string `json:"tenant_id"`
    SellerID  string `json:"seller_id"`
    Quantity  int    `json:"quantity"`
}
```

---

## Integration Testing

- Integration tests are in `controllers/*_test.go`.
- Tests use Gin’s test mode and can be run with:
  ```sh
  go test ./controllers/...
  ```
- For real DB/Redis integration, ensure those services are running (e.g., via Docker Compose).

---

## Notes

- All endpoints expect and return JSON.
- Use query parameters for filtering (e.g., `tenant_id`, `seller_id`, `sku_code`).
- Redis is used for caching; cache is invalidated on write operations.

---

If you need more details or want to add usage examples, let me know!
