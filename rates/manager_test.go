package rates

import (
	"fmt"
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"testing"
)

type MockProvider struct {
	slug          string
	updater       func() error
	ratesProvider func() *dto.Rates
}

func NewMockProvider(slug string, updater func() error, ratesProvider func() *dto.Rates) *MockProvider {
	return &MockProvider{
		slug:          slug,
		updater:       updater,
		ratesProvider: ratesProvider,
	}
}

func (p *MockProvider) Update(ch chan<- error) {
	ch <- p.updater()
}

func (p *MockProvider) GetSlug() string {
	return p.slug
}

func (p *MockProvider) GetRates() *dto.Rates {
	return p.ratesProvider()
}

func TestNewManager(t *testing.T) {
	m := NewManager()

	if len(m.providers) != 0 {
		t.Errorf("providers should be empty array, got %d elements", len(m.providers))
	}
}

func TestManager_AddProvider(t *testing.T) {
	m := NewManager()

	slug := "slug"
	p := NewMockProvider(slug, func() error {
		return nil
	}, func() *dto.Rates {
		return nil
	})
	err := m.AddProvider(p)
	if err != nil {
		t.Errorf("adding of provider first time should be permitted, got error: %s", err)
	}
	err = m.AddProvider(p)
	if err == nil || err.Error() != "provider with slag slug already exist, overwriting is restricted" {
		t.Error("adding of provider with the same slag should be restricted")
	}
	if len(m.providers) != 1 {
		t.Errorf("providers should contain added provider, got %d elements", len(m.providers))
	}
	_, ok := m.providers[slug]
	if !ok {
		t.Errorf("cannot retrieve provider by slug %s", slug)
	}
}

func TestManager_GetRates(t *testing.T) {
	m := NewManager()

	_, err := m.GetRates("wrongSlug")
	if err != nil && err.Error() != "provider with slag wrongSlug does not exist" {
		t.Errorf("provider does not exist, got error %s", err.Error())
	}

	slug := "slug"
	base := "USD"
	expectedRates := dto.NewRates(base)
	p := NewMockProvider(slug, func() error {
		return nil
	}, func() *dto.Rates {
		return expectedRates
	})
	err = m.AddProvider(p)
	if err != nil {
		t.Errorf("adding of provider first time should be permitted, got error: %s", err)
	}

	rates, _ := m.GetRates(slug)
	if rates != expectedRates {
		t.Errorf("rates should equal, expected: %s, got: %s", expectedRates, rates)
	}
}

func TestManager_Update(t *testing.T) {
	m := NewManager()

	slug := "slug"
	expectedErr := fmt.Errorf("expected error")
	p := NewMockProvider(slug, func() error {
		return expectedErr
	}, func() *dto.Rates {
		return nil
	})
	err := m.AddProvider(p)
	if err != nil {
		t.Errorf("adding of provider first time should be permitted, got error: %s", err)
	}

	err = m.Update()
	if err == nil {
		t.Error("expected error should happen, got nothing")
	}
	if err != nil && err != expectedErr {
		t.Errorf("expected error should happen, got error: %s", err)
	}

	expectedErr = nil
	err = m.Update()
	if err != nil {
		t.Errorf("no errors should happen, got error: %s", err)
	}
}
