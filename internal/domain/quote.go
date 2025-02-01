package domain

import (
	"math/rand"
	"time"
)

type Quote string

func (r Quote) Bytes() []byte {
	return []byte(r)
}

type QuoteProvider interface {
	Provide() Quote
}

type RandomQuoteProvider struct {
	phrases []Quote
}

func NewRandomPhraseGrantProvider(quotes []Quote) *RandomQuoteProvider {
	return &RandomQuoteProvider{
		quotes,
	}
}

func (r *RandomQuoteProvider) Provide() Quote {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(r.phrases))
	return r.phrases[index]
}
