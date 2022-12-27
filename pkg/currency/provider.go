package currency

import (
	"context"
	"net/http"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

type Provider interface {
	GetCurrencyExchangeFactor(ctx context.Context, base, second string) (float64, error)
}

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetProvidersPool will return map of pointers to providers ready to be used
func GetProvidersPool(
	client HttpDoer,
	providerConfigs map[string]config.ProviderConfig,
) map[string]Provider {
	providers := map[string]Provider{}

	openExchange, exists := providerConfigs[NameOpenExchangeRates]
	if exists {
		providers[NameOpenExchangeRates] = NewOpenExchangeProvider(client, openExchange.AppKey)
	}

	currencyConverter, exists := providerConfigs[NameCurrencyConverter]
	if exists {
		providers[NameCurrencyConverter] = NewCurrencyConverterProvider(client, currencyConverter.AppKey)
	}

	return providers
}
