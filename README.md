# Omniful Assignment - Microservices Architecture

A comprehensive microservices-based e-commerce platform built with Go, featuring Inventory Management System (IMS) and Order Management System (OMS) services. This project demonstrates modern microservices architecture with event-driven design, multi-database support, and cloud-native patterns.

## 🏗️ Architecture Overview

This project implements a microservices architecture with the following components:

- **IMS (Inventory Management System):** Manages product SKUs, warehouse hubs, and inventory levels
- **OMS (Order Management System):** Handles bulk order processing, CSV uploads, and order lifecycle management
- **Supporting Infrastructure:** PostgreSQL, MongoDB, Redis, Kafka, S3, SQS, and LocalStack

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.20+ (for local development)
- Git

### Running the Complete System

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd Omniful-Assignment
   ```

2. **Start all services:**
   ```bash
   docker-compose up -d
   ```

3. **Verify services are running:**
   ```bash
   docker-compose ps
   ```

## 📋 Service Details

### IMS Service (Inventory Management System)
- **Port:** 8084
- **Database:** PostgreSQL
- **Cache:** Redis
- **Features:** SKU management, hub management, inventory tracking

**Quick Test:**
```bash
curl http://localhost:8084/health
```

### OMS Service (Order Management System)
- **Port:** 8086
- **Database:** MongoDB
- **Message Queue:** SQS
- **Event Streaming:** Kafka
- **File Storage:** S3
- **Features:** Bulk order processing, CSV uploads, order management

**Quick Test:**
```bash
curl http://localhost:8086/health
```

## 🗄️ Infrastructure Services

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5434 | Primary database for IMS |
| MongoDB | 27018 | Document database for OMS |
| Redis | 6380 | Caching layer |
| Kafka | 9092 | Event streaming |
| Zookeeper | 2181 | Kafka coordination |
| LocalStack | 4566 | AWS service emulation |

## 📁 Project Structure

```
Omniful-Assignment/
├── docker-compose.yml          # Infrastructure orchestration
├── README.md                   # This file
├── ims/                        # Inventory Management System
│   ├── README.md              # IMS documentation
│   ├── main.go                # IMS service entry point
│   ├── controllers/           # HTTP request handlers
│   ├── models/                # Data models
│   ├── db/                    # Database layer
│   ├── migrations/            # Database migrations
│   ├── routes/                # API routes
│   └── redisclient/           # Redis integration
├── oms/                        # Order Management System
│   ├── README.md              # OMS documentation
│   ├── main.go                # OMS service entry point
│   ├── controllers/           # HTTP request handlers
│   ├── models/                # Data models
│   ├── database/              # Database layer
│   ├── routes/                # API routes
│   ├── utils/                 # External service clients
│   ├── middleware/            # HTTP middleware
│   └── localization/          # Internationalization
├── go_commons/                 # Shared Go utilities
├── add_inventory_columns.sql   # Sample data setup
└── *.csv                       # Test data files
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6380

# MongoDB Configuration
MONGODB_URI=mongodb://myuser:mypassword@localhost:27018/mydb?authSource=admin
MONGODB_DB_NAME=mydb

# AWS/LocalStack Configuration
AWS_REGION=us-east-1
AWS_S3_ENDPOINT=http://localhost:4566
AWS_SQS_ENDPOINT=http://localhost:4566
S3_BUCKET_NAME=order-csv-bucket
CREATE_BULK_ORDER_QUEUE_NAME=CreateBulkOrder

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_ORDER_TOPIC=order.created

# Service URLs
IMS_BASE_URL=http://localhost:8084
```

## 📊 API Documentation

### IMS API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/hub/` | Create hub |
| `GET` | `/hub/` | List hubs |
| `POST` | `/sku/` | Create SKU |
| `GET` | `/sku/` | List SKUs |
| `POST` | `/inventory/` | Create inventory |
| `GET` | `/inventory/` | Get inventory |
| `PUT` | `/inventory/upsert` | Upsert inventory |

### OMS API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/orders/upload` | Upload CSV for bulk orders |
| `GET` | `/orders/` | List orders |
| `GET` | `/orders/:id` | Get order by ID |
| `PUT` | `/orders/:id/status` | Update order status |

## 🧪 Testing

### Sample Data

The project includes several CSV files for testing:

- `valid_test_with_real_data.csv` - Valid order data
- `valid_combinations.csv` - Valid SKU-location combinations
- `test_working_inventory.csv` - Test inventory data
- `invalid_empty.csv` - Empty file for error testing
- `invalid_missing_columns.csv` - Missing required columns

### Quick API Tests

**Test IMS Health:**
```bash
curl http://localhost:8084/health
```

**Test OMS Health:**
```bash
curl http://localhost:8086/health
```

**Upload CSV to OMS:**
```bash
curl -X POST http://localhost:8086/orders/upload \
  -F "file=@valid_test_with_real_data.csv"
```

**Get Inventory from IMS:**
```bash
curl "http://localhost:8084/inventory/?sku=SKU23&location=HU001"
```

## 🔄 Data Flow

### Order Processing Flow

1. **CSV Upload:** Client uploads CSV file to OMS
2. **File Storage:** OMS uploads file to S3
3. **Message Queue:** OMS publishes S3 path to SQS
4. **Background Processing:** SQS consumer processes CSV
5. **Inventory Validation:** OMS validates with IMS
6. **Order Creation:** Orders stored in MongoDB
7. **Event Publishing:** Order events published to Kafka
8. **Inventory Update:** Kafka consumer updates inventory via IMS

### Inventory Management Flow

1. **SKU Creation:** Create product SKUs with metadata
2. **Hub Management:** Define warehouse locations
3. **Inventory Tracking:** Track stock levels by SKU and location
4. **Real-time Updates:** Atomic upsert operations
5. **Caching:** Redis cache for performance optimization

## 🛠️ Development

### Local Development Setup

1. **Start infrastructure:**
   ```bash
   docker-compose up -d postgres redis mongodb kafka zookeeper localstack
   ```

2. **Run IMS service:**
   ```bash
   cd ims
   go run main.go
   ```

3. **Run OMS service:**
   ```bash
   cd oms
   go run main.go
   ```

### Database Migrations

**IMS Migrations:**
```bash
cd ims
go run db/migration.go
```

**Sample Data Setup:**
```bash
docker exec -i postgres_db psql -U postgres -d postgres < add_inventory_columns.sql
```

## 📈 Monitoring and Observability

- **Health Checks:** `/health` endpoints on both services
- **Structured Logging:** Comprehensive logging with correlation IDs
- **Metrics:** Built-in metrics for performance monitoring
- **Error Tracking:** Detailed error logging and handling

## 🔒 Security Features

- **Input Validation:** Comprehensive validation across all endpoints
- **File Type Validation:** Strict CSV file validation
- **Authentication Ready:** Middleware prepared for authentication
- **Rate Limiting:** Built-in rate limiting capabilities

## 🚀 Deployment

### Production Considerations

1. **Environment Variables:** Configure production environment variables
2. **Database Security:** Use production-grade databases with proper security
3. **Message Queues:** Replace LocalStack with actual AWS services
4. **Monitoring:** Implement proper monitoring and alerting
5. **Load Balancing:** Add load balancers for high availability

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose -f docker-compose.yml up -d

# Scale services if needed
docker-compose up -d --scale ims=2 --scale oms=2
```

## 🤝 Contributing

1. Follow Go coding standards and conventions
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Ensure proper error handling and logging
5. Test with various data formats and edge cases

## 📚 Additional Resources

- [IMS Documentation](./ims/README.md) - Detailed IMS service documentation
- [OMS Documentation](./oms/README.md) - Detailed OMS service documentation
- [Docker Compose](./docker-compose.yml) - Infrastructure configuration
- [Sample Data](./*.csv) - Test data files

## 📄 License

This project is part of the Omniful Assignment and follows the project's licensing terms.

## 🆘 Support

For issues and questions:
1. Check the service-specific README files
2. Review the API documentation
3. Check the health endpoints for service status
4. Review logs for detailed error information

