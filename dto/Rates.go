//go:generate easyjson -all
package dto

import "sync"

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
