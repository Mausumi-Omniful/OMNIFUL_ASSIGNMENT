# IMS (Inventory Management System)

A Go-based microservice for inventory management using the `omniful/go_commons` library.

## Features

- **SKU Management:** CRUD operations for product SKUs, with filtering by tenant, seller, and SKU code.
- **Hub Management:** CRUD operations for warehouse locations, with filtering by tenant and seller.
- **Inventory Management:** CRUD operations for inventory levels, atomic upsert, and default quantity 0 for missing entries.
- **Multi-tenant Architecture:** Support for `tenant_id` and `seller_id`.
- **Database Migrations:** Automatic schema management.
- **Redis Caching:** For SKU and hub validation.
- **Docker Support:** Containerized deployment.

## Project Structure

```
Omniful-Assignment/
  ├─ docker-compose.yml
  ├─ ims/
  │   ├─ controllers/
  │   │   ├─ hub_controller.go
  │   │   ├─ inventory_controller.go
  │   │   └─ sku_controller.go
  │   ├─ db/
  │   │   ├─ connection.go
  │   │   └─ migration.go
  │   ├─ main.go
  │   ├─ go.mod
  │   ├─ go.sum
  │   ├─ migrations/
  │   │   ├─ 20240618_create_hub_sku_inventory.up.sql
  │   │   └─ 20240618_create_hub_sku_inventory.down.sql
  │   ├─ models/
  │   │   ├─ hub.go
  │   │   ├─ inventory.go
  │   │   └─ sku.go
  │   ├─ redisclient/
  │   │   └─ redis_client.go
  │   └─ routes/
  │       └─ routes.go
  └─ README.md
```

## Prerequisites

- Docker and Docker Compose
- Go 1.20+ (for local development)

## Running with Docker

1. **Clone the repository and navigate to the project directory:**
   ```bash
   git clone <repo-url>
   cd Omniful-Assignment
   ```

2. **Start all services using Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **The IMS service will be available at:**
   - **URL:** http://localhost:8083

## Running Locally (without Docker)

1. **Install dependencies:**
   ```bash
   cd ims
   go mod download
   ```

2. **Set up environment variables:**
   - Copy `.env.example` to `.env` (if available) or create a `.env` file with your configuration.

3. **Run the application:**
   ```bash
   go run main.go
   ```

## API Endpoints

### Hub Management
- `POST /hub/` - Create a new hub
- `GET /hub/` - Get all hubs (with optional filters: `tenant_id`, `seller_id`)
- `GET /hub/:id` - Get hub by ID
- `PUT /hub/:id` - Update hub
- `DELETE /hub/:id` - Delete hub

### SKU Management
- `POST /sku/` - Create a new SKU
- `GET /sku/` - Get all SKUs (with optional filters: `tenant_id`, `seller_id`, `sku_code`)
- `PUT /sku/:id` - Update SKU
- `DELETE /sku/:id` - Delete SKU

### Inventory Management
- `POST /inventory/` - Create a new inventory item
- `GET /inventory/` - Get inventory by SKU and location (query params: `sku`, `location`)
- `GET /inventory/:id` - Get inventory by ID
- `PUT /inventory/:id` - Update inventory
- `PUT /inventory/upsert` - Atomic upsert (insert or update by SKU and location)
- `DELETE /inventory/:id` - Delete inventory

## Environment Variables

The application uses the following environment variables (set in `.env`):

- `DB_HOST` - PostgreSQL host (default: postgresql)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `REDIS_HOST` - Redis host (default: redis)
- `REDIS_PORT` - Redis port (default: 6379)

## Database Migrations

Migrations are automatically run when the application starts.  
To manually run migrations:

```bash
go run db/migration.go
# To rollback:
go run db/migration.go rollback
```

## Testing the API

You can use Postman or curl to test the endpoints.  
Example for checking inventory:
```
GET http://localhost:8083/inventory/?sku=TESTSKU&location=TESTLOC
```

## License

MIT 