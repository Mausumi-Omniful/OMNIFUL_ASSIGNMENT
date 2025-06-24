package responses

import "time"

// FetchRatesRes Being Used in Exchange Rate Service as well
type FetchRatesRes struct {
	BaseCurrency string             `json:"base_currency"`
	Rates        map[string]float64 `json:"rates"`
	LastUpdated  time.Time          `json:"last_updated"`
}
