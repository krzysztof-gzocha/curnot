package aggregator

import (
	"bytes"
	"fmt"

	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/notifier"
)

type RateAggregatorInterface interface {
	Aggregate(newRate *Rate) error
}

type RateAggregator struct {
	lastRate        map[string]*Rate
	notifierChain   notifier.NotifierChain
	currencyConfigs []config.CurrencyConfig
}

func NewRateAggregator(
	notifierChain *notifier.NotifierChain,
	currencyConfig []config.CurrencyConfig,
) *RateAggregator {
	return &RateAggregator{
		notifierChain:   *notifierChain,
		currencyConfigs: currencyConfig,
		lastRate:        make(map[string]*Rate),
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

		err := ra.notify(newRate)

		if err != nil {
			return err
		}
	}

	ra.lastRate[newRate.String()] = newRate

	return nil
}

func (ra *RateAggregator) notify(newRate *Rate) error {
	msg := bytes.NewBufferString(fmt.Sprintf(
		"1 %s = %.4f %s",
		newRate.From,
		newRate.Rate,
		newRate.To,
	))

	last, exists := ra.lastRate[newRate.String()]
	if !exists || last.Rate == 0 || newRate.Rate == last.Rate {
		return ra.notifierChain.Notify(msg.String())
	}

	// Add " (+12%)" to the notification
	compare := newRate.Rate / last.Rate
	if newRate.Rate > last.Rate {
		msg.WriteString(" (+")
	} else {
		msg.WriteString(" (-")
	}
	msg.WriteString(fmt.Sprintf("%.2f%%)", compare*100))

	return ra.notifierChain.Notify(msg.String())
}
