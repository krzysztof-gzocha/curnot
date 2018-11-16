package currency

import (
	"net/http"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

type ProviderInterface interface {
	GetCurrencyExchangeFactor(base, second string) (float64, error)
}

type HttpClientInterface interface {
	Get(url string) (resp *http.Response, err error)
}

// GetProvidersPool will return map of pointers to providers ready to be used
func GetProvidersPool(
	client HttpClientInterface,
	providerConfigs map[string]config.ProviderConfig,
) map[string]ProviderInterface {
	providers := map[string]ProviderInterface{}

	openExchange, exists := providerConfigs[NameOpenExchangeRates]
	if exists {
		providers[NameOpenExchangeRates] = NewOpenExchangeProvider(client, openExchange.AppKey)
	}

	_, exists = providerConfigs[NameCurrencyConverter]
	if exists {
		providers[NameCurrencyConverter] = NewCurrencyConverterProvider(client)
	}

	return providers
}
