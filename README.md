# IMS (Inventory Management System)

A Go-based microservice for inventory management using the `omniful/go_commons` library.

## Features

- **SKU Management**: CRUD operations for product SKUs
- **Hub Management**: CRUD operations for warehouse locations
- **Inventory Management**: CRUD operations for inventory levels
- **Multi-tenant Architecture**: Support for tenant_id and seller_id
- **Database Migrations**: Automatic schema management
- **Redis Caching**: Performance optimization
- **Docker Support**: Containerized deployment

## Prerequisites

- Docker and Docker Compose
- Go 1.24.3+ (for local development)

## Quick Start with Docker

1. **Clone the repository and navigate to the project directory**

2. **Start all services using Docker Compose:**
   ```bash
   docker-compose --env-file docker.env up -d
   ```

3. **The IMS service will be available at:**
   - **URL**: http://localhost:8083
   - **Health Check**: http://localhost:8083/health

## API Endpoints

### SKU Management
- `POST /sku/` - Create a new SKU
- `GET /sku/` - Get all SKUs (with optional filters)
- `GET /sku/:id` - Get SKU by ID
- `PUT /sku/:id` - Update SKU
- `DELETE /sku/:id` - Delete SKU

### Hub Management
- `POST /hub/` - Create a new hub
- `GET /hub/` - Get all hubs (with optional filters)
- `GET /hub/:id` - Get hub by ID
- `PUT /hub/:id` - Update hub
- `DELETE /hub/:id` - Delete hub

### Inventory Management
- `POST /inventory/` - Create a new inventory item
- `GET /inventory/` - Get all inventory items (with optional filters)
- `GET /inventory/:id` - Get inventory by ID
- `PUT /inventory/:id` - Update inventory
- `DELETE /inventory/:id` - Delete inventory

## Query Parameters

### SKU Filters
- `tenant_id` - Filter by tenant ID
- `seller_id` - Filter by seller ID
- `sku_code` - Filter by SKU code

### Hub Filters
- `tenant_id` - Filter by tenant ID
- `seller_id` - Filter by seller ID

### Inventory Filters
- `tenant_id` - Filter by tenant ID
- `seller_id` - Filter by seller ID
- `sku_id` - Filter by SKU ID
- `hub_id` - Filter by hub ID

## Environment Variables

The application uses the following environment variables:

### Database Configuration
- `DB_HOST` - PostgreSQL host (default: postgresql)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name

### Redis Configuration
- `REDIS_HOST` - Redis host (default: redis)
- `REDIS_PORT` - Redis port (default: 6379)

## Services

The Docker Compose setup includes:

- **IMS Service** (Port 8083) - Main application
- **PostgreSQL** (Port 5432) - Primary database
- **Redis** (Port 6379) - Caching layer
- **MongoDB** (Port 27017) - Document database (configured but not used)
- **Kafka** (Port 9092) - Event streaming (configured but not used)
- **Zookeeper** (Port 2181) - Kafka coordination

## Local Development

1. **Install dependencies:**
   ```bash
   cd ims
   go mod download
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

## Database Migrations

Migrations are automatically run when the application starts. To manually run migrations:

```bash
# Run migrations up
go run db/migration.go

# Rollback migrations
go run db/migration.go rollback
```

## Architecture

- **Models**: Data structures and GORM tags
- **Controllers**: Business logic and HTTP handlers
- **Routes**: API endpoint definitions
- **Database**: Connection management and migrations
- **Redis**: Caching layer
- **Docker**: Containerization and orchestration

## Dependencies

- **Gin**: Web framework
- **GORM**: ORM for database operations
- **omniful/go_commons**: Shared utilities and database connections
- **godotenv**: Environment variable management 