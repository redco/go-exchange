//go:generate easyjson -all
package dto

import (
	"errors"
	"fmt"
	"sync"
)

//easyjson:json
type Rates struct {
	mu    sync.RWMutex
	Items map[string]float32 `json:"rates"`
	Base  string             `json:"base"`
}

func NewRates(base string) *Rates {
	return &Rates{
		Items: make(map[string]float32, 0),
		Base:  base,
	}
}

func (r *Rates) Add(currency string, rate float32) {
	r.mu.Lock()
	r.Items[currency] = rate
	r.mu.Unlock()
}

func (r *Rates) String() string {
	return fmt.Sprintf("{ Base: %s, Items: %s }", r.Base, r.Items)
}

func (r *Rates) AddFromMap(jsonMap map[string]map[string]float32) error {
	r.mu.Lock()
	for currency, rate := range jsonMap {
		var value float32
		var ok bool
		if value, ok = rate[r.Base]; !ok {
			return errors.New("cannot fetch rate by base currency")
		}
		r.Items[currency] = 1 / value
	}
	r.mu.Unlock()
	return nil
}
