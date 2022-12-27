package main

import (
	"context"
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

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("Version: %v\nCommit %v\nBuilt at %v\n\n", version, commit, date)

	ctx := context.Background()
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
	notifierChain := notifier.NewChain(httpClient, cfg.Notifiers)

	agg := aggregator.NewRateAggregator(notifierChain, cfg.Currencies)
	tickerChecker := checker.NewTickerChecker(
		cfg.Interval,
		checker.NewChecker(cfg.Currencies, providersPool, agg),
	)

	fmt.Printf("Starting checking every %s..", cfg.Interval)
	err = tickerChecker.Check(ctx)
	fmt.Println(err.Error())
}
