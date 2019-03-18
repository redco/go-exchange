package rates

import (
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"io/ioutil"
	"log"
	"net/http"
)

// Provider represents any provider which can be plugged into the Manager
// to have multi rates provider system.
type Provider interface {
	// Update updates rates from remote server and keeps it locally.
	Update(ch chan<- error)
	// GetSlug returns slug of particular rates provider, which should be
	// unique among all other providers registered within Manager.
	GetSlug() string
	// GetRates returns cached local copy of provider rates.
	GetRates() *dto.Rates
}

type HttpFetcher struct{}

// Fetcher is an interface representing the ability to fetch
// remote url address and receive content from it.
type Fetcher interface {
	// Fetch fetches url and returns it's data.
	// It returns []bytes body and error if any.
	Fetch(url string) ([]byte, error)
}

// Fetch fetches url data via http.Get.
// It returns []bytes body and error if any.
func (f *HttpFetcher) Fetch(url string) (bodyBytes []byte, err error) {
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
