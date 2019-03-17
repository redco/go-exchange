package rates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"log"
)

const urlsChunk = 3

type CryptoCompare struct {
	slug    string
	urls    []string
	rates   dto.Rates
	fetcher Fetcher
}

func NewCryptoCompare(base string, currencies []string, fetcher Fetcher) *CryptoCompare {
	urls, err := getUrlsFromCurrencies(base, currencies)
	if err != nil {
		log.Fatal(err)
	}
	return &CryptoCompare{
		slug:    "cryptocompare",
		urls:    urls,
		rates:   *dto.NewRates(base),
		fetcher: fetcher,
	}
}

func (p *CryptoCompare) Update(ch chan<- error) {
	urlsCount := len(p.urls)
	done := make(chan error, urlsCount)

	for i := 0; i < urlsCount; i++ {
		url := p.urls[i]
		go func() {
			bodyBytes, err := p.fetcher.fetch(url)
			if err != nil {
				done <- err
				return
			}
			jsonMap := make(map[string]map[string]float32)
			err = json.Unmarshal(bodyBytes, &jsonMap)
			if err != nil {
				done <- err
				return
			}
			err = p.rates.AddFromMap(jsonMap)
			if err != nil {
				done <- err
				return
			}
			log.Println(jsonMap)
			log.Printf("Fetched result: %s", bodyBytes)
			done <- nil
		}()
	}

	for i := 0; i < urlsCount; i++ {
		err := <-done
		if err != nil {
			ch <- err
			return
		}
	}
	ch <- nil
}

func (p *CryptoCompare) GetSlug() string {
	return p.slug
}

func (p *CryptoCompare) GetRates() *dto.Rates {
	return &p.rates
}

func getUrlsFromCurrencies(base string, currencies []string) ([]string, error) {
	var buf bytes.Buffer
	var urls []string
	for i := 1; i <= len(currencies); i++ {
		buf.WriteString(currencies[i-1])
		if i%int(urlsChunk) == 0 || i == len(currencies) {
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
