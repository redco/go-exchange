package pkg

import (
	"encoding/json"
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"io/ioutil"
	"log"
	"net/http"
)

type RatesManger struct {
	rates dto.Rates
	urls  []string
}

func NewRatesManager(base string, urls []string) *RatesManger {
	return &RatesManger{
		rates: *dto.NewRates(base),
		urls:  urls,
	}
}

func (m *RatesManger) Update() {
	urlsCount := len(m.urls)
	done := make(chan error, urlsCount)

	for i := 0; i < urlsCount; i++ {
		url := m.urls[i]
		go func() {
			bodyBytes, err := fetch(url)
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
			err = m.rates.AddFromMap(jsonMap)
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
			log.Fatal(err)
		}
	}
}

func (m *RatesManger) Get() dto.Rates {
	return m.rates
}

func fetch(url string) (bodyBytes []byte, err error) {
	log.Printf("fetching url: %s", url)
	resp, err := http.Get(url)
	defer func() {
		errClose := resp.Body.Close()
		if err == nil {
			err = errClose
		}
	}()
	if err != nil {
		return
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}
