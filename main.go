package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/jinzhu/configor"
	"github.com/krzysztof-gzocha/curnot/pkg/aggregator"
	"github.com/krzysztof-gzocha/curnot/pkg/checker"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
	"github.com/krzysztof-gzocha/curnot/pkg/currency"
	"github.com/krzysztof-gzocha/curnot/pkg/notifier"
)

func main() {
	var cfgFileName string
	flag.StringVar(&cfgFileName, "cfg", path.Join(os.Getenv("HOME"), "curnot.yaml"), "Config file")
	flag.Parse()

	if cfgFileName == "" {
		fmt.Println("Please specify config file path by --cfg")
		os.Exit(1)
		return
	}
	f, err := os.OpenFile(cfgFileName, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("No such file: %s\n", cfgFileName)
		os.Exit(1)
		return
	}
	_ = f.Close()

	ctx := context.Background()
	cfg := config.Config{}
	err = configor.New(&configor.Config{ENVPrefix: "CURNOT"}).Load(&cfg, cfgFileName)
	if err != nil {
		fmt.Printf("Could not load %s: %+v\n", cfgFileName, err)
		os.Exit(1)
		return
	}

	httpTimeout := cfg.Timeout
	if cfg.Interval < cfg.Timeout {
		httpTimeout = cfg.Interval
	}

	httpClient := &http.Client{Timeout: httpTimeout}
	providersPool := currency.GetProvidersPool(httpClient, cfg.Providers)
	notifierChain := notifier.NewChain(httpClient, cfg.Notifiers)

	agg := aggregator.NewRateAggregator(notifierChain, cfg.Currencies)
	tickerChecker := checker.NewTickerChecker(
		cfg.Interval,
		checker.NewChecker(cfg.Currencies, providersPool, agg),
	)

	fmt.Printf("Starting checking every %s..\n\n", cfg.Interval)

	tickerChecker.StartChecking(ctx)
}
