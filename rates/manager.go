package rates

import (
	"fmt"
	"github.com/maZahaca/go-exchange-rates-service/dto"
)

// Manager contains multiple rates providers and providers
// ability to Update and GetRates of each of it.
type Manager struct {
	providers map[string]Provider
}

// NewManager creates new Manager and returns reference for it.
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

// AddProvider registers new provider in manager for further
// using via Update and GetRates methods.
// It returns error if any happened.
func (m *Manager) AddProvider(p Provider) error {
	_, ok := m.providers[p.GetSlug()]
	if ok {
		return fmt.Errorf("provider with slag %s already exist, overwriting is restricted", p.GetSlug())
	}
	m.providers[p.GetSlug()] = p
	return nil
}

// Update calls Update method for every registered provider.
// It returns error if any happened.
func (m *Manager) Update() error {
	count := len(m.providers)
	done := make(chan error, count)

	for _, p := range m.providers {
		go p.Update(done)
	}

	var err error
	for i := 0; i < count; i++ {
		err = <-done
	}
	if err != nil {
		return err
	}
	return nil
}

// GetRates returns rates from specified provider.
// It returns Rates and error if any happened.
func (m *Manager) GetRates(slug string) (*dto.Rates, error) {
	p, ok := m.providers[slug]
	if !ok {
		return nil, fmt.Errorf("provider with slag %s does not exist", slug)
	}
	return p.GetRates(), nil
}
