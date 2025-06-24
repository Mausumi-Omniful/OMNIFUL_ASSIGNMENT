package exchange_rate

import (
	"context"
	commonErr "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/exchange_rate/responses"
	"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
)

// Paths
const (
	getExchangeRatesByCurrency = "/internal/api/v1/exchange_rates/{baseCurrency}"
)

// Path Params
const (
	baseCurrencyParam = "baseCurrency"
)

type ExchangeRateClient interface {
	FetchRates(ctx context.Context, baseCurrency string) (responses.FetchRatesRes, commonErr.CustomError)
}

type exchangeRateInterserviceClient struct {
	*interservice_client.Client
}

func newExchangeRateInterserviceClient(intersvcClient *interservice_client.Client) *exchangeRateInterserviceClient {
	return &exchangeRateInterserviceClient{
		intersvcClient,
	}
}

func (client *exchangeRateInterserviceClient) FetchRates(ctx context.Context, baseCurrency string) (responses.FetchRatesRes, commonErr.CustomError) {
	res := responses.FetchRatesRes{}

	_, intersvcErr := client.Execute(ctx, http.APIGet, &http.Request{
		Url:        getExchangeRatesByCurrency,
		PathParams: map[string]string{baseCurrencyParam: baseCurrency},
	}, &res)
	if intersvcErr != nil {
		return responses.FetchRatesRes{}, commonErr.NewCustomError(commonErr.BadRequestError, intersvcErr.Message)
	}

	return res, commonErr.CustomError{}
}
