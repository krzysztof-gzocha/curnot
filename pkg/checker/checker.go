package checker

import (
	"fmt"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/currency"
	"github.com/pkg/errors"
)

type CheckerInterface interface {
	Check() error
}

type Checker struct {
	currencies []config.CurrencyConfig
	providers  map[string]currency.ProviderInterface
	aggregator aggregator.RateAggregatorInterface
}

func NewChecker(
	currencies []config.CurrencyConfig,
	providers map[string]currency.ProviderInterface,
	aggregator aggregator.RateAggregatorInterface,
) *Checker {
	return &Checker{
		currencies: currencies,
		providers:  providers,
		aggregator: aggregator,
	}
}

func (c *Checker) Check() error {
	for _, currencyConfig := range c.currencies {
		provider, exists := c.providers[currencyConfig.ProviderName]
		if !exists {
			fmt.Printf("Provider %s for %s%s does not exist", currencyConfig.ProviderName, currencyConfig.From, currencyConfig.To)
			continue
		}

		currencyRate, err := provider.GetCurrencyExchangeFactor(currencyConfig.From, currencyConfig.To)
		if err != nil {
			return errors.Wrapf(err, "Could not fetch currency rate from provider '%s'", currencyConfig.ProviderName)
		}

		c.aggregator.Aggregate(&aggregator.Rate{
			From: currencyConfig.From,
			To:   currencyConfig.To,
			Rate: currencyRate,
		})
	}

	return nil
}
