package rates

import (
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"io/ioutil"
	"log"
	"net/http"
)

type Provider interface {
	Update(ch chan<- error)
	GetSlug() string
	GetRates() *dto.Rates
}

type HttpFetcher struct{}

type Fetcher interface {
	fetch(url string) ([]byte, error)
}

func (f *HttpFetcher) fetch(url string) (bodyBytes []byte, err error) {
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
