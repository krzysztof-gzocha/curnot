package config

import (
	"time"
)

type Config struct {
	Interval   time.Duration
	Providers  map[string]ProviderConfig
	Currencies []CurrencyConfig
	Notifiers  map[string]NotifierConfig
}

type NotifierConfig struct {
	EmailReceiverParameters EmailReceiverParameters   `yaml:"receiver"`
	ConnectionParameters    EmailConnectionParameters `yaml:"connection_parameters"`
}

type EmailConnectionParameters struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailReceiverParameters struct {
	Email string
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
