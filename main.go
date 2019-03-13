package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maZahaca/go-exchange-rates-service/pkg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	baseCurrency = flag.String("base", "USD", "base currency to get rates base on it")
)

func main() {
	flag.Parse()

	urlsChunk, err := strconv.ParseInt(getEnv("URL_CHUNKS", "3"), 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	urls, err := getUrlsFromArgs(flag.Args(), urlsChunk, *baseCurrency)
	if err != nil {
		log.Fatal(err)
	}

	manager := pkg.NewRatesManager(*baseCurrency, urls)
	manager.Update()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			manager.Update()
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rawResponse, err := manager.Get().MarshalJSON()
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

func getUrlsFromArgs(args []string, urlsChunk int64, base string) ([]string, error) {
	var buf bytes.Buffer
	var urls []string
	for i := 1; i <= len(args); i++ {
		buf.WriteString(args[i-1])
		if i%int(urlsChunk) == 0 || i == len(args) {
			urls = append(
				urls,
				fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=%s",
					buf.String(), base),
			)
			buf.Reset()
			continue
		}
		buf.WriteString(",")
	}
	return urls, nil
}
