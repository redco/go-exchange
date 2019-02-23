# go-exchange-rates-service

This service keeps exchange rates from [CryptoCompare.com](https://cryptocompare.com).
It refreshes it periodically and serve it as JSON REST API.
The service allows to specify `base` currency and other currencies which rates need to be kept.

## Run
```bash
go run main.go --base USD -- GBP EUR CAD RUB CHF BTC ETH ETC
```