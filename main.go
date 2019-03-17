package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maZahaca/go-exchange-rates-service/rates"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	baseCurrency = flag.String("base", "USD", "base currency to get rates base on it")
)

func main() {
	flag.Parse()

	fetcher := new(rates.HttpFetcher)
	manager := rates.NewManager()
	provider := rates.NewCryptoCompare(*baseCurrency, flag.Args(), fetcher)
	err := manager.AddProvider(provider)
	if err != nil {
		log.Fatal(err)
	}

	err = manager.Update()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			err = manager.Update()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		pRates, err := manager.GetRates("cryptocompare")
		if err != nil {
			log.Fatal(err)
		}
		rawResponse, err := pRates.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		if _, err = writer.Write(rawResponse); err != nil {
			log.Fatal(err)
		}
	})

	port := getEnv("PORT", "8000")
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func getEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}
