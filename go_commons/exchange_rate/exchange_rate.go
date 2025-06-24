package exchange_rate

import (
	"context"
	"fmt"
	interserviceClient "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"github.com/omniful/go_commons/shutdown"
	"sync"
	"time"
)

const (
	DefaultBaseCurrency    = "USD"
	DefaultUpdateFrequency = time.Hour
)

type ExchangeRateServiceConfig struct {
	baseCurrency    string
	updateFrequency time.Duration
}

func WithBaseCurrency(baseCurrency string) func(*ExchangeRateServiceConfig) {
	return func(c *ExchangeRateServiceConfig) {
		c.baseCurrency = baseCurrency
	}
}

func WithUpdateFrequency(duration time.Duration) func(*ExchangeRateServiceConfig) {
	return func(c *ExchangeRateServiceConfig) {
		c.updateFrequency = duration
	}
}

// exchangeRate represents a single exchange rate between two currencies
type exchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

// ExchangeRateService Service handles exchange rate operations and caching
type ExchangeRateService struct {
	client     ExchangeRateClient
	rates      map[string]exchangeRate // FromCurrency -> ToCurrency -> Rate
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	config     ExchangeRateServiceConfig
}

// NewService creates a new exchange service with the given client and base currency
func NewExchangeRateServiceWithInterserviceClient(ctx context.Context,
	exchangeClient *interserviceClient.Client,
	opts ...func(*ExchangeRateServiceConfig)) (*ExchangeRateService, error) {
	exchangeRateClient := newExchangeRateInterserviceClient(exchangeClient)

	return NewExchangeRateService(ctx, exchangeRateClient, opts...)
}

// NewExchangeRateService creates a new exchange service with the given client and config, starts update rate poller
func NewExchangeRateService(ctx context.Context,
	exchangeClient ExchangeRateClient,
	opts ...func(*ExchangeRateServiceConfig)) (*ExchangeRateService, error) {
	ctx, cancel := context.WithCancel(ctx)

	config := ExchangeRateServiceConfig{
		baseCurrency:    DefaultBaseCurrency,
		updateFrequency: DefaultUpdateFrequency,
	}

	for _, opt := range opts {
		opt(&config)
	}

	exchangeRateSvc := &ExchangeRateService{
		client:     exchangeClient,
		rates:      make(map[string]exchangeRate),
		ctx:        ctx,
		cancelFunc: cancel,
		config:     config,
	}

	err := exchangeRateSvc.start()
	if err != nil {
		return exchangeRateSvc, err
	}

	return exchangeRateSvc, nil
}

// Start begins periodic updates of exchange rates
func (s *ExchangeRateService) start() error {
	if s == nil {
		return fmt.Errorf("ExchangeRateService is nil")
	}

	// Start periodic updates
	go s.periodicUpdate()

	shutdown.RegisterShutdownCallback("exchange-rate-service", s)

	// Fetch initial rates
	if err := s.updateRates(); err != nil {
		return fmt.Errorf("initial rate fetch failed: %s", err.Error())
	}

	return nil
}

// Close terminates the periodic updates
func (s *ExchangeRateService) Close() error {
	if s == nil {
		return fmt.Errorf("ExchangeRateService is nil")
	}

	s.cancelFunc()

	return nil
}

func (s *ExchangeRateService) BaseCurrency() string {
	if s == nil {
		return ""
	}

	return s.config.baseCurrency
}

// convertFromBase converts an amount from base currency to target currency
func (s *ExchangeRateService) convertFromBase(toCurrency string, amount float64) (float64, error) {
	if s == nil {
		return 0, fmt.Errorf("ExchangeRateService is nil")
	}

	if toCurrency == s.BaseCurrency() {
		return amount, nil
	}

	exchangeRate, ok := s.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("no exchange rate found for %s to %s", s.BaseCurrency(), toCurrency)
	}

	return amount * exchangeRate.Rate, nil
}

// convertToBase converts an amount from target currency to baseCurrency currency
func (s *ExchangeRateService) convertToBase(fromCurrency string, amount float64) (float64, error) {
	if s == nil {
		return 0, fmt.Errorf("ExchangeRateService is nil")
	}

	if fromCurrency == s.BaseCurrency() {
		return amount, nil
	}

	exchangeRate, ok := s.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("no exchange rate found for %s to %s", s.BaseCurrency(), fromCurrency)
	}

	return amount / exchangeRate.Rate, nil
}

// Convert converts an amount from one currency to another using base currency as intermediate
func (s *ExchangeRateService) Convert(fromCurrency, toCurrency string, amount float64) (float64, error) {
	if s == nil {
		return 0, fmt.Errorf("ExchangeRateService is nil")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Handle direct conversion if either currency is the base currency
	if fromCurrency == s.BaseCurrency() {
		return s.convertFromBase(toCurrency, amount)
	} else if toCurrency == s.BaseCurrency() {
		return s.convertToBase(fromCurrency, amount)
	}

	// Convert through base currency for cross-currency conversion
	baseAmount, err := s.convertToBase(fromCurrency, amount)
	if err != nil {
		return 0, err
	}

	finalAmount, err := s.convertFromBase(toCurrency, baseAmount)
	if err != nil {
		return 0, fmt.Errorf("conversion path error %s -> %s -> %s: %s",
			fromCurrency, s.BaseCurrency(), toCurrency, err.Error())
	}

	return finalAmount, nil
}

// periodicUpdate handles the periodic refresh of exchange rates
func (s *ExchangeRateService) periodicUpdate() {
	if s == nil {
		return
	}

	ticker := time.NewTicker(s.config.updateFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			if err := s.updateRates(); err != nil {
				newrelic.NoticeError(s.ctx, fmt.Errorf("exchange rates updates failed"))

				continue
			}
		}
	}
}

// updateRates fetches and stores new exchange rates
func (s *ExchangeRateService) updateRates() error {
	if s == nil {
		return fmt.Errorf("ExchangeRateService is nil")
	}

	res, cusErr := s.client.FetchRates(s.ctx, s.config.baseCurrency)
	if cusErr.Exists() {
		log.Errorf("failed to fetch rates: %s", cusErr.Error())

		return cusErr.ToError()
	}

	newRates := make(map[string]exchangeRate)
	for toCurrency, rate := range res.Rates {
		newRates[toCurrency] = exchangeRate{
			FromCurrency: res.BaseCurrency,
			ToCurrency:   toCurrency,
			Rate:         rate,
			LastUpdated:  res.LastUpdated,
		}
	}

	s.mu.Lock()
	s.rates = newRates
	defer s.mu.Unlock()

	return nil
}
