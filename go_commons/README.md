# **Go Commons**

[Code Coverage](https://raw.githubusercontent.com/omniful/go_commons/badges/.badges/master/coverage.svg)

Go Commons is a comprehensive collection of Golang packages providing essential functionality for building robust microservices and distributed systems. This repository contains various utility packages that can be used across different Omniful services.

## Core Packages

### Message Queue and Worker Management
- [**SQS**](./sqs/): AWS SQS integration package with support for message publishing, consumption, batching, and automatic compression.
- [**Worker**](./worker/): A robust framework for managing background tasks and concurrent job processing with support for multiple listener types (HTTP, Kafka, SQS).
- [**Kafka**](./kafka/): Kafka integration package for message streaming and event processing.

### Data Storage and Caching
- [**Redis**](./redis/): Redis client utilities and helper functions.
- [**Redis Cache**](./redis_cache/): Caching implementation using Redis.
- [**S3**](./s3/): AWS S3 integration utilities for object storage.
- [**DB**](./db/): Database utilities and helpers.

### HTTP and API
- [**HTTP**](./http/): HTTP server and client utilities.
- [**HTTPClient**](./httpclient/): HTTP client implementation with advanced features.
- [**JWT**](./jwt/): JSON Web Token handling utilities.
- [**Interservice-Client**](./interservice-client/): Client for inter-service communication.

### Monitoring and Observability
- [**Monitoring**](./monitoring/): System monitoring utilities.
- [**NewRelic**](./newrelic/): NewRelic integration package.
- [**Health**](./health/): Health check implementations.
- [**Log**](./log/): Logging utilities.

### Data Processing
- [**JSON**](./json/): JSON handling utilities.
- [**CSV**](./csv/): CSV file processing utilities.
- [**Compression**](./compression/): Data compression utilities.
- [**File Utilities**](./file_utilities/): File handling and processing utilities.

### Configuration and Environment
- [**Config**](./config/): Configuration management utilities.
- [**Env**](./env/): Environment variable handling.
- [**Constants**](./constants/): Common constants used across services.

### Utilities
- [**Error**](./error/): Error handling utilities.
- [**I18n**](./i18n/): Internationalization support.
- [**Pagination**](./pagination/): Pagination implementation.
- [**Permissions**](./permissions/): Permission management utilities.
- [**Pool**](./pool/): Resource pooling utilities.
- [**PubSub**](./pubsub/): Publisher/Subscriber pattern implementation.
- [**Ratelimiter**](./ratelimiter/): Rate limiting implementation.
- [**Set**](./set/): Set data structure implementation.
- [**Shutdown**](./shutdown/): Graceful shutdown utilities.
- [**Mobile Number**](./mobile_number/): Mobile number handling utilities.
- [**Exchange Rate**](./exchange_rate/): Currency exchange rate utilities.
- [**Currency Converter**](./currencyconverter/): Currency conversion utilities.
- [**Nullable**](./nullable/): Nullable type implementations.
- [**Rules**](./rules/): Business rules engine.
- [**Runtime**](./runtime/): Runtime utilities.
- [**Weight**](./weight/): Weight calculation and conversion utilities.
- [**DChannel**](./dchannel/): Distributed channel implementation.
- [**DMutex**](./dmutex/): Distributed mutex implementation.

## Installation

```bash
go get github.com/omniful/go_commons
```

## Development Setup

### 1. Set the GOPRIVATE Environment Variable
To ensure Go treats the repository as private and allows access, run the following command in your terminal:

```bash
go env -w GOPRIVATE="github.com/omniful/"
```

### 2. Configure Git Settings
Modify or create the `~/.gitconfig` file with the following details, replacing `{name}` with your Git username and `{email}` with your Git email:

```ini
[user]
    name = {name}
    email = {email}

[url "ssh://git@github.com/"]
    insteadOf = https://github.com/
```

This configuration instructs Git to use SSH instead of HTTPS for GitHub URLs, which is necessary for accessing private repositories.

### 3. Configure SSH for GitHub
Make sure your SSH key is set up and added to your GitHub account. This is required for accessing private repositories.

## Contributing

Please read through our contributing guidelines before making any contributions.

## License

This project is licensed under the terms of our license agreement.
