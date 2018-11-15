package currency

import "net/http"

type ProviderInterface interface {
	GetCurrencyExchangeFactor(base, second string) (float32, error)
}

type HttpClientInterface interface {
	Get(url string) (resp *http.Response, err error)
}
