package config

import (
	"time"
)

type Config struct {
	Interval   time.Duration
	Providers  map[string]ProviderConfig
	Currencies []CurrencyConfig
}

type ProviderConfig struct {
	AppKey string `yaml:"app_key"`
}

type Alert struct {
	AnyChange bool `yaml:"any_change"`
	Below     float64
	Above     float64
}

type CurrencyConfig struct {
	From         string
	To           string
	ProviderName string `yaml:"provider_name"`
	Alert        Alert
}
