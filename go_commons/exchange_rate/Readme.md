# Exchange Rate Package

## Overview
The `exchange_rate` package provides a robust solution for handling currency exchange rates in Go applications. It offers a service-based approach with features like automatic rate updates, caching, and flexible currency conversions. The package is designed to work with an interservice client for fetching real-time exchange rates.

## Features
- Real-time exchange rate fetching
- Automatic periodic rate updates
- In-memory caching of exchange rates
- Cross-currency conversion support
- Configurable base currency
- Thread-safe operations
- Graceful shutdown handling

## Installation
```go
import "github.com/omniful/go_commons/exchange_rate"
```

## Configuration
The exchange rate service can be configured with the following options:

```go
// Default configuration
const (
    DefaultBaseCurrency    = "USD"
    DefaultUpdateFrequency = time.Hour
)

// Custom configuration
service, err := exchange_rate.NewExchangeRateServiceWithInterserviceClient(
    ctx,
    interserviceClient,
    exchange_rate.WithBaseCurrency("EUR"),
    exchange_rate.WithUpdateFrequency(30 * time.Minute),
)
```

## Usage Examples

### 1. Creating a New Service
```go
package main

import (
    "context"
    "github.com/omniful/go_commons/exchange_rate"
    interserviceClient "github.com/omniful/go_commons/interservice-client"
)

func main() {
    ctx := context.Background()
    
    // Initialize with interservice client
    client := // your interservice client initialization
    service, err := exchange_rate.NewExchangeRateServiceWithInterserviceClient(ctx, client)
    if err != nil {
        panic(err)
    }
    defer service.Close()
}
```

### 2. Converting Currencies
```go
// Convert 100 USD to EUR
amount := 100.0
converted, err := service.Convert("USD", "EUR", amount)
if err != nil {
    log.Printf("Error converting currency: %v", err)
    return
}
fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
```

## API Reference

### Types

#### ExchangeRateService
The main service struct that handles all exchange rate operations.

### Methods

#### NewExchangeRateServiceWithInterserviceClient
```go
func NewExchangeRateServiceWithInterserviceClient(
    ctx context.Context,
    exchangeClient *interserviceClient.Client,
    opts ...func(*ExchangeRateServiceConfig),
) (*ExchangeRateService, error)
```
Creates a new exchange rate service with an interservice client.

#### Convert
```go
func (s *ExchangeRateService) Convert(
    fromCurrency string,
    toCurrency string,
    amount float64,
) (float64, error)
```
Converts an amount from one currency to another using the base currency as an intermediate if needed.

#### Close
```go
func (s *ExchangeRateService) Close() error
```
Gracefully shuts down the service and stops periodic updates.

### Configuration Options

#### WithBaseCurrency
```go
func WithBaseCurrency(baseCurrency string) func(*ExchangeRateServiceConfig)
```
Sets the base currency for the service (default: "USD").

#### WithUpdateFrequency
```go
func WithUpdateFrequency(duration time.Duration) func(*ExchangeRateServiceConfig)
```
Sets how often the exchange rates should be updated (default: 1 hour).

## Notes
- The service automatically starts updating exchange rates upon initialization
- All operations are thread-safe
- The service uses the configured base currency for intermediate conversions
- Exchange rates are cached in memory and updated periodically
- The service implements graceful shutdown through the `shutdown` package

## Error Handling
The service returns appropriate errors in the following cases:
- Invalid currency codes
- Missing exchange rates
- Service initialization failures
- Rate fetching failures

## Best Practices
1. Always close the service when done using `defer service.Close()`
2. Handle errors appropriately in production code
3. Choose an appropriate update frequency based on your needs
4. Consider your base currency carefully as all cross-currency conversions go through it

## How It Works

### Periodic Rate Updates
The exchange rate service implements an automatic polling mechanism to keep exchange rates up-to-date:

1. **Initialization**: When the service starts, it immediately fetches initial rates
2. **Background Polling**: A goroutine runs in the background that:
   - Polls at configured intervals (default: 1 hour)
   - Updates the in-memory cache with new rates
   - Handles errors gracefully without disrupting the service
3. **Thread-safe Cache**: All rates are stored in an in-memory map protected by a RWMutex
4. **Graceful Shutdown**: The polling goroutine is properly cleaned up when the service is closed

```go
// Internal structure of cached rates
type exchangeRate struct {
    FromCurrency string
    ToCurrency   string
    Rate         float64
    LastUpdated  time.Time
}
```

### Complete Example with Interservice Client
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/omniful/go_commons/exchange_rate"
    "github.com/omniful/go_commons/interservice-client"
    "github.com/omniful/go_commons/http"
)

func main() {
    ctx := context.Background()

    // Initialize the interservice client
    client, err := interservice_client.NewClient(
        "exchange-rate-service",  // service name
        &http.Config{
            BaseURL:    "https://your-exchange-rate-service.com",
            Timeout:    5 * time.Second,
            RetryCount: 3,
        },
    )
    if err != nil {
        log.Fatalf("Failed to initialize interservice client: %v", err)
    }

    // Create exchange rate service with custom configuration
    service, err := exchange_rate.NewExchangeRateServiceWithInterserviceClient(
        ctx,
        client,
        exchange_rate.WithBaseCurrency("USD"),
        exchange_rate.WithUpdateFrequency(30 * time.Minute),
    )
    if err != nil {
        log.Fatalf("Failed to initialize exchange rate service: %v", err)
    }
    defer service.Close()

    // The service is now:
    // 1. Automatically fetching rates every 30 minutes
    // 2. Caching them in memory
    // 3. Ready for currency conversions

    // Example conversion
    amount := 1000.0
    fromCurrency := "USD"
    toCurrency := "EUR"

    converted, err := service.Convert(fromCurrency, toCurrency, amount)
    if err != nil {
        log.Printf("Error converting currency: %v", err)
        return
    }

    fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, converted, toCurrency)

    // Service will automatically clean up when closed
}
```

### Internal Rate Updates
The service maintains rates through several key mechanisms:

1. **Initial Load**
```go
// On service start, rates are immediately loaded
err := exchangeRateSvc.updateRates()
if err != nil {
    return fmt.Errorf("initial rate fetch failed: %s", err.Error())
}
```

2. **Periodic Updates**
```go
// Background goroutine for updates
func (s *ExchangeRateService) periodicUpdate() {
    ticker := time.NewTicker(s.config.updateFrequency)
    defer ticker.Stop()

    for {
        select {
        case <-s.ctx.Done():
            return
        case <-ticker.C:
            if err := s.updateRates(); err != nil {
                // Error handling with retry logic
                continue
            }
        }
    }
}
```

3. **Thread-safe Access**
```go
// Example of thread-safe rate access
func (s *ExchangeRateService) Convert(fromCurrency, toCurrency string, amount float64) (float64, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // ... conversion logic ...
}
```
