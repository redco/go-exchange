# go-exchange
[![CircleCI (all branches)](https://img.shields.io/circleci/project/github/redco/go-exchange.svg)](https://circleci.com/gh/redco/go-exchange)
[![Codecov](https://img.shields.io/codecov/c/github/redco/go-exchange.svg)](https://codecov.io/gh/redco/go-exchange)
[![GoDoc](https://godoc.org/github.com/redco/go-exchange?status.svg)](https://godoc.org/github.com/redco/go-exchange)
[![Go Report Card](https://goreportcard.com/badge/github.com/redco/go-exchange)](https://goreportcard.com/report/github.com/redco/go-exchange)


This service keeps exchange rates from [CryptoCompare.com](https://cryptocompare.com).
It refreshes it periodically and serve it as JSON REST API.
The service allows to specify `base` currency and other currencies which rates need to be kept.

## Run
```bash
go run main.go --base USD -- GBP EUR CAD RUB CHF BTC ETH ETC
```
