// Copyright 2019 RedCode. All rights reserved.
// This is an example library demonstrating different
// Go language techniques and practices.

/*
Package rates implements Manager and Rates Provider.
This allows your to add multiple rates providers and update it via manager.
All the rates are stored in the memory and can be accessed via GetRates method.
Let's have a look on example of using this package:
	fetcher := new(rates.HttpFetcher)
	manager := rates.NewManager()
	provider := rates.NewCryptoCompare("USD", []string{"GBP", "EUR", "BTC"}, fetcher)
	err := manager.AddProvider(provider)
	if err != nil {
		log.Fatal(err)
	}

	err = manager.Update()
	if err != nil {
		log.Fatal(err)
	}
	rates := manager.GetRates("cryptocompare")

Here we create Manager and register with it CryptoCompare provider, every time
we call Update method the rates are fetched and renewed in the memory.
*/

package rates
