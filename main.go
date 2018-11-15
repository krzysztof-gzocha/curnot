package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/configor"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/checker"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/currency"
	"github.com/krzysztof-gzocha/curnot/pkg/notifier"
)

const configFile = "config.yml"
const timeout = time.Second * 10
const tickInterval = time.Hour

func main() {
	cfg := config.Config{}
	err := configor.Load(&cfg, configFile)
	if err != nil {
		fmt.Printf("Could not load %s: %s\n", configFile, err.Error())
		return
	}

	httpClient := &http.Client{
		Timeout: timeout,
	}

	providersPool := currency.GetProvidersPool(httpClient, cfg.Providers)
	desktopNotifier := notifier.NewDesktop()
	ticker := time.NewTicker(tickInterval)
	agg := aggregator.NewRateAggregator(desktopNotifier, cfg.Currencies)
	tickerChecker := checker.NewTickerChecker(
		ticker,
		checker.NewChecker(cfg.Currencies, providersPool, agg),
	)

	err = tickerChecker.Check()
	fmt.Println(err.Error())
}
