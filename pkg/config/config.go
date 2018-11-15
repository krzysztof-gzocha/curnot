package config

type Config struct {
	Providers  map[string]ProviderConfig
	Currencies []CurrencyConfig
}

type ProviderConfig struct {
	AppKey string
}

type CurrencyConfig struct {
	From         string
	To           string
	ProviderName string
}
