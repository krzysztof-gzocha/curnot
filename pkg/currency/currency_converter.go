package currency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const NameCurrencyConverter = "currencyConverter"

type CurrencyConverterProvider struct {
	client HttpClientInterface
}

func NewCurrencyConverterProvider(client HttpClientInterface) *CurrencyConverterProvider {
	return &CurrencyConverterProvider{
		client: client,
	}
}

func (c *CurrencyConverterProvider) GetCurrencyExchangeFactor(
	base,
	second string,
) (float64, error) {
	response, err := c.client.Get(c.buildUrl(base, second))
	if err != nil {
		return 0, err
	}

	content, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, errors.Errorf(
			"Response status: %d, content: %s",
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
		return 0, errors.Errorf("Response is not formatted correctly: %s", string(content))
	}

	rate, exists := currencyResponse["val"]
	if !exists {
		return 0, errors.Errorf("Response is not formatted correctly: %s", string(content))
	}

	return rate, nil
}

func (c *CurrencyConverterProvider) getCurrencyName(base, second string) string {
	return fmt.Sprintf("%s_%s", base, second)
}

// https://free.currencyconverterapi.com/api/v5/convert?q=USD_PLN&compact=y
func (c *CurrencyConverterProvider) buildUrl(base, second string) string {
	values := &url.Values{}
	values.Set("q", c.getCurrencyName(base, second))
	values.Set("compact", "y")

	uri := &url.URL{}
	uri.Scheme = "https"
	uri.Host = "free.currencyconverterapi.com"
	uri.Path = "/api/v5/convert"
	uri.RawQuery = values.Encode()

	return uri.String()
}
