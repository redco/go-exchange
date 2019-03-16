package rates

import (
	"fmt"
	"github.com/maZahaca/go-exchange-rates-service/dto"
	"log"
)

type Manger struct {
	providers map[string]Provider
}

func NewManager() *Manger {
	return &Manger{
		providers: make(map[string]Provider),
	}
}

func (m *Manger) AddProvider(p Provider) {
	m.providers[p.GetSlug()] = p
}

func (m *Manger) Update() {
	count := len(m.providers)
	done := make(chan error, count)

	for _, p := range m.providers {
		go p.Update(done)
	}

	for i := 0; i < count; i++ {
		err := <-done
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *Manger) GetRates(slug string) (*dto.Rates, error) {
	p, ok := m.providers[slug]
	if !ok {
		return nil, fmt.Errorf("provider with slag %s does not exist", slug)
	}
	return p.GetRates(), nil
}
