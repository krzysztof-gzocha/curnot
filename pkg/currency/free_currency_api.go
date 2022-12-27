package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const NameFreeCurrencyApi = "freeCurrencyApi"

type FreeCurrencyApi struct {
	client HttpDoer
	apiKey string
}

type freeCurrencyApiResponse struct {
	Rates map[string]float64 `json:"data"`
}

func NewFreeCurrencyApi(client HttpDoer, apiKey string) *FreeCurrencyApi {
	return &FreeCurrencyApi{client: client, apiKey: apiKey}
}

func (f *FreeCurrencyApi) GetCurrencyExchangeFactor(ctx context.Context, base, second string) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, f.buildUrl(base, second), nil)
	if err != nil {
		return 0, err
	}

	resp, err := f.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}

	respStruct := freeCurrencyApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return 0, err
	}

	rate, exists := respStruct.Rates[second]
	if !exists {
		return 0, fmt.Errorf("rate %s/%s was missing from response", base, second)
	}

	return rate, nil
}

func (f *FreeCurrencyApi) buildUrl(base, second string) string {
	values := &url.Values{}
	values.Set("apikey", f.apiKey)
	values.Set("base", base)
	values.Set("currencies", second)

	address := &url.URL{}
	address.Scheme = "https"
	address.Host = "api.freecurrencyapi.com"
	address.Path = "/v1/latest"
	address.RawQuery = values.Encode()

	return address.String()
}
