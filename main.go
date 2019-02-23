package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	baseCurrency = flag.String("base", "USD", "base currency to get rates base on it")
	urlsChunk    = 3
)

func main() {
	flag.Parse()

	var buf bytes.Buffer
	var urls []string
	args := flag.Args()
	for i := 1; i <= len(args); i++ {
		buf.WriteString(args[i-1])
		if i%urlsChunk == 0 || i == len(args) {
			log.Println("i=", i)
			urls = append(
				urls,
				fmt.Sprintf("https://min-api.cryptocompare.com/data/pricemulti?fsyms=%s&tsyms=%s",
					buf.String(), *baseCurrency),
			)
			buf.Reset()
			continue
		}
		buf.WriteString(",")
	}

	rates := dto.NewRates(*baseCurrency)

	updateRates(urls, *rates)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			updateRates(urls, *rates)
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rawResponse, err := rates.MarshalJSON()
		if err != nil {
			log.Fatal(err)
		}
		writer.Write(rawResponse)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func updateRates(urls []string, rates dto.Rates) {
	done := make(chan string, len(urls))

	for i := 0; i < len(urls); i++ {
		go fetch(urls[i], done)
	}

	for i := 0; i < len(urls); i++ {
		res := <-done
		jsonMap := make(map[string]map[string]float32)
		err := json.Unmarshal([]byte(res), &jsonMap)
		if err != nil {
			panic(err)
		}
		for currency, rate := range jsonMap {
			rates.Add(currency, 1/rate[*baseCurrency])
		}
		log.Println(rates)
		log.Printf("Fetched result: %s", res)
	}
}

func fetch(url string, done chan string) {
	log.Printf("fetching url: %s", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Cannot fetch url: %+q", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cannot read response: %+q", err)
	}
	done <- string(bodyBytes)
}
