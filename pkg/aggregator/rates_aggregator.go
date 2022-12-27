package aggregator

import (
	"context"
	"fmt"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

type RateAggregator interface {
	Aggregate(ctx context.Context, newRate *Rate) error
}

type MultipleRatesAggregator struct {
	lastRate        map[string]*Rate
	notifier        Notifier
	currencyConfigs []config.CurrencyConfig
}

func NewRateAggregator(
	notifier Notifier,
	currencyConfig []config.CurrencyConfig,
) *MultipleRatesAggregator {
	return &MultipleRatesAggregator{
		notifier:        notifier,
		currencyConfigs: currencyConfig,
		lastRate:        make(map[string]*Rate),
	}
}

func (ra *MultipleRatesAggregator) Aggregate(ctx context.Context, newRate *Rate) error {
	for _, currencyConfig := range ra.currencyConfigs {
		if !newRate.supports(currencyConfig) {
			continue
		}

		if !newRate.shouldNotify(currencyConfig.Alert) {
			fmt.Println("Shouldn't notify, skipping..")
			continue
		}

		err := ra.notify(ctx, newRate)
		if err != nil {
			return err
		}
	}

	ra.lastRate[newRate.String()] = newRate

	return nil
}

func (ra *MultipleRatesAggregator) notify(ctx context.Context, newRate *Rate) error {
	rc := RateChange{New: newRate}
	if last, exists := ra.lastRate[newRate.String()]; exists {
		rc.Old = last
	}

	return ra.notifier.Notify(ctx, rc)
}
