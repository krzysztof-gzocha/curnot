package config

type Config struct {
	Providers  map[string]ProviderConfig
	Currencies []CurrencyConfig
}

type ProviderConfig struct {
	AppKey string
}

type Alert struct {
	AnyChange bool
	Below     float64
	Above     float64
}

type CurrencyConfig struct {
	From         string
	To           string
	ProviderName string
	Alert        Alert
}
