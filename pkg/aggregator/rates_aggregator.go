package aggregator

import (
	"fmt"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/notifier"
)

type RateAggregatorInterface interface {
	Aggregate(newRate *Rate) error
}

type RateAggregator struct {
	notifier        notifier.NotifierInterface
	currencyConfigs []config.CurrencyConfig
}

func NewRateAggregator(
	notifier notifier.NotifierInterface,
	currencyConfig []config.CurrencyConfig,
) *RateAggregator {
	return &RateAggregator{
		notifier:        notifier,
		currencyConfigs: currencyConfig,
	}
}

func (ra *RateAggregator) Aggregate(newRate *Rate) error {
	for _, currencyConfig := range ra.currencyConfigs {
		if !newRate.supports(currencyConfig) {
			continue
		}

		if !newRate.shouldNotify(currencyConfig.Alert) {
			continue
		}

		err := ra.notifier.Notify(fmt.Sprintf(
			"1 %s = %.2f %s",
			newRate.From,
			newRate.Rate,
			newRate.To,
		))

		if err != nil {
			return err
		}
	}

	return nil
}
