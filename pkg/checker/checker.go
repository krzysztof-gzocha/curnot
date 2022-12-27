package checker

import (
	"context"
	"fmt"

	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/currency"
	"github.com/pkg/errors"
)

type Checker interface {
	Check(ctx context.Context) error
}

type ProvidersChecker struct {
	currencies []config.CurrencyConfig
	providers  map[string]currency.Provider
	aggregator aggregator.RateAggregator
}

func NewChecker(
	currencies []config.CurrencyConfig,
	providers map[string]currency.Provider,
	aggregator aggregator.RateAggregator,
) *ProvidersChecker {
	return &ProvidersChecker{
		currencies: currencies,
		providers:  providers,
		aggregator: aggregator,
	}
}

func (c *ProvidersChecker) Check(ctx context.Context) error {
	for _, currencyConfig := range c.currencies {
		provider, exists := c.providers[currencyConfig.ProviderName]
		if !exists {
			fmt.Printf("Provider %s for %s%s does not exist\n", currencyConfig.ProviderName, currencyConfig.From, currencyConfig.To)
			continue
		}

		fmt.Printf(
			"Using provider '%s' to check for %s%s\n",
			currencyConfig.ProviderName,
			currencyConfig.From,
			currencyConfig.To,
		)

		currencyRate, err := provider.GetCurrencyExchangeFactor(ctx, currencyConfig.From, currencyConfig.To)
		if err != nil {
			return errors.Wrapf(err, "Could not fetch currency rate from provider '%s'", currencyConfig.ProviderName)
		}

		fmt.Printf(
			"Provider '%s' returned 1 %s = %.4f %s\n",
			currencyConfig.ProviderName,
			currencyConfig.From,
			currencyRate,
			currencyConfig.To,
		)

		err = c.aggregator.Aggregate(ctx, &aggregator.Rate{
			From: currencyConfig.From,
			To:   currencyConfig.To,
			Rate: currencyRate,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
