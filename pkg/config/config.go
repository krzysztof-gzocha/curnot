package config

import (
	"time"
)

type Config struct {
	Interval   time.Duration
	Providers  map[string]ProviderConfig
	Currencies []CurrencyConfig
	Notifiers  map[string]NotifierConfig // @todo improve, so each notifier will have only its own cfg
}

type NotifierConfig struct {
	EmailReceiverParameters EmailReceiverParameters   `yaml:"receiver"`
	ConnectionParameters    EmailConnectionParameters `yaml:"connection_parameters"`
	HttpParameters          HttpParams                `yaml:"http_parameters"`
}

type HttpParams struct {
	Method                   string `yaml:"method"`
	Path                     string `yaml:"path"`
	AcceptedResponseStatuses []int  `yaml:"accepted_response_statuses"`
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
