package rates

import (
	"errors"
	"testing"
)

var (
	baseCurrency = "USD"
	currencies   = []string{"GBP"}
)

type MockFetcher struct {
	fetcher func(url string) ([]byte, error)
}

func NewMockFetcher(fetcher func(url string) ([]byte, error)) *MockFetcher {
	return &MockFetcher{
		fetcher: fetcher,
	}
}

func (mf *MockFetcher) Fetch(url string) ([]byte, error) {
	return mf.fetcher(url)
}

func getCryptoCompare(bytes []byte, err error) *CryptoCompare {
	f := NewMockFetcher(func(url string) ([]byte, error) {
		return bytes, err
	})

	return NewCryptoCompare(baseCurrency, currencies, f)
}

func TestNewCryptoCompare(t *testing.T) {
	p := getCryptoCompare([]byte(`{"GBP":{"USD":1.34}}`), nil)

	if len(p.urls) != 1 {
		t.Error("urls are not initialized")
	}

	expectedUrl := "https://min-api.cryptocompare.com/data/pricemulti?fsyms=GBP&tsyms=USD"
	if p.urls[0] != expectedUrl {
		t.Errorf("urls are initialized incorrect, expected: %s, got: %s", expectedUrl, p.urls[0])
	}

	expectedSlug := "cryptocompare"
	if p.GetSlug() != expectedSlug {
		t.Errorf("slug is initialized incorrect, expected: %s, got: %s", expectedSlug, p.GetSlug())
	}
}

func TestCryptoCompare_Update(t *testing.T) {
	p := getCryptoCompare([]byte(`{"GBP":{"USD":1.34}}`), nil)

	done := make(chan error, 1)
	p.Update(done)
	err := <-done

	if err != nil {
		t.Errorf("error is happened, expected: nil, got: %s", err)
	}

	errExpected := errors.New("expected error")
	p = getCryptoCompare(nil, errExpected)
	done = make(chan error, 1)
	p.Update(done)
	err = <-done

	if err != errExpected {
		t.Errorf("error is happened, expected: %s, got: %s", errExpected, err)
	}

	p = getCryptoCompare([]byte(``), nil)
	done = make(chan error, 1)
	p.Update(done)
	err = <-done

	if err.Error() != "unexpected end of JSON input" {
		t.Errorf("error is happened, expected: %s, got: %s", errExpected, err)
	}
}

func TestCryptoCompare_GetRates(t *testing.T) {
	p := getCryptoCompare([]byte(`{"GBP":{"USD":1.34}}`), nil)

	rates := *p.GetRates()
	if rates.Base != "USD" {
		t.Error("wong base currency, expected USD")
	}
	if len(rates.Items) != 0 {
		t.Error("rates should be empty before update")
	}

	done := make(chan error, 1)
	p.Update(done)
	err := <-done
	if err != nil {
		t.Errorf("error should be empty, got %s", err)
	}
	rates = *p.GetRates()
	if len(rates.Items) != 1 {
		t.Error("rates should have one value")
	}
	expectedRate := float32(0.74626863)
	if rates.Items["GBP"] != expectedRate {
		t.Errorf("error is happened, expected: %g, got: %g", expectedRate, rates.Items["GBP"])
	}
}
