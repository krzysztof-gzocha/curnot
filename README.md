# Currency Notifier
Small golang program capable of periodically checking different currencies exchange rates with different providers.

# Implemented providers
- currencyConverter - https://free.currencyconverterapi.com
- openExchangeRates - https://openExchangeRates.org

## Have a look at releases to download compiled version

# Config example

```yaml
# Configure interval of notifications. Default is 30 minutes
# Examples:
# 60s => 60 seconds
# 30m => 30 minutes
# 2h => 2 hours
interval: 30m

# Providers section is configuring all possible providers that you can use.
# I recommend to use currencyConverter by default, but if you need openExchangeRates.org
# you just need to uncomment the configuration and paste your application key
providers:
  currencyConverter: ~
#  openExchangeRates:
#    app_key: "<PASTE YOUR KEY HERE>"

# Currencies section is configuring all the currencies you would like to track.
# Feel free to add more, but 1 currency rate will result in 1 notification.
currencies:
  - from: USD
    to: PLN
    provider_name: currencyConverter
    alert:
      below: 3.75
      above: 3.85
      #any_change: true
```

# Project is under development
