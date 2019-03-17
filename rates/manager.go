package rates

import (
	"fmt"
	"github.com/maZahaca/go-exchange-rates-service/dto"
)

type Manager struct {
	providers map[string]Provider
}

func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
	}
}

func (m *Manager) AddProvider(p Provider) error {
	_, ok := m.providers[p.GetSlug()]
	if ok {
		return fmt.Errorf("provider with slag %s already exist, overwriting is restricted", p.GetSlug())
	}
	m.providers[p.GetSlug()] = p
	return nil
}

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

func (m *Manager) GetRates(slug string) (*dto.Rates, error) {
	p, ok := m.providers[slug]
	if !ok {
		return nil, fmt.Errorf("provider with slag %s does not exist", slug)
	}
	return p.GetRates(), nil
}
