package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const NameCurrencyConverter = "currencyConverter"

type CurrencyConverterProvider struct {
	client HttpDoer
	apiKey string
}

func NewCurrencyConverterProvider(client HttpDoer, apiKey string) *CurrencyConverterProvider {
	return &CurrencyConverterProvider{
		client: client,
		apiKey: apiKey,
	}
}

func (c *CurrencyConverterProvider) GetCurrencyExchangeFactor(
	ctx context.Context,
	base,
	second string,
) (float64, error) {
	req, err := http.NewRequest(http.MethodGet, c.buildUrl(c.apiKey, base, second), nil)
	if err != nil {
		return 0, err
	}
	response, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}

	content, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf(
			"response status: %d, content: %s",
			response.StatusCode,
			string(content),
		)
	}

	responseBody := map[string]map[string]float64{}
	err = json.Unmarshal(content, &responseBody)
	if err != nil {
		return 0, err
	}

	currencyResponse, exists := responseBody[c.getCurrencyName(base, second)]
	if !exists {
		return 0, fmt.Errorf("response is not formatted correctly: %s", string(content))
	}

	rate, exists := currencyResponse["val"]
	if !exists {
		return 0, fmt.Errorf("response is not formatted correctly: %s", string(content))
	}

	return rate, nil
}

func (c *CurrencyConverterProvider) getCurrencyName(base, second string) string {
	return fmt.Sprintf("%s_%s", base, second)
}

// https://free.currencyconverterapi.com/api/v5/convert?q=USD_PLN&compact=y
func (c *CurrencyConverterProvider) buildUrl(apiKey, base, second string) string {
	values := &url.Values{}
	values.Set("q", c.getCurrencyName(base, second))
	values.Set("compact", "y")
	values.Set("apiKey", apiKey)

	uri := &url.URL{}
	uri.Scheme = "https"
	uri.Host = "free.currencyconverterapi.com"
	uri.Path = "/api/v5/convert"
	uri.RawQuery = values.Encode()

	return uri.String()
}
