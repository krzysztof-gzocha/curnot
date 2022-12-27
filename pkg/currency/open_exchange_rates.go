package currency

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const NameOpenExchangeRates = "openExchangeRates"

type OpenExchangeProvider struct {
	client HttpDoer
	apiKey string
}

type openExchangeResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func NewOpenExchangeProvider(client HttpDoer, apiKey string) *OpenExchangeProvider {
	return &OpenExchangeProvider{client: client, apiKey: apiKey}
}

func (cp *OpenExchangeProvider) GetCurrencyExchangeFactor(ctx context.Context, base, second string) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, cp.buildUrl(cp.apiKey, base, second), nil)
	if err != nil {
		return 0, err
	}
	response, err := cp.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}

	content, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return 0, errors.Wrapf(err, "Could not read response body for %s/%s", base, second)
	}

	if response.StatusCode != http.StatusOK {
		return 0, errors.Errorf("Response status: %d. Content: '%s'", response.StatusCode, string(content))
	}

	responseBody := openExchangeResponse{}
	err = json.Unmarshal(content, &responseBody)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not unmarshal response body: %s", content)
	}

	rate, exists := responseBody.Rates[second]
	if !exists {
		return 0, errors.Errorf("Rate for %s/%s was missing in the response: %s", base, second, string(content))
	}

	return rate, nil
}

func (cp *OpenExchangeProvider) buildUrl(apiKey, base, second string) string {
	values := &url.Values{}
	address := &url.URL{}
	values.Set("app_id", apiKey)
	values.Set("base", base)
	values.Set("symbols", second)

	address.Scheme = "https"
	address.Host = "openexchangerates.org"
	address.Path = "/api/latest.json"
	address.RawQuery = values.Encode()

	return address.String()
}
