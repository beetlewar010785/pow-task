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

type WordOfWisdomQuoteProvider struct {
	phrases []Quote
}

func NewWordOfWisdomQuoteProvider() *WordOfWisdomQuoteProvider {
	return &WordOfWisdomQuoteProvider{
		[]Quote{
			"Cease to be idle; cease to be unclean; cease to find fault one with another.",
			"A man is saved no faster than he gains knowledge.",
			"Our thoughts determine our actions, our actions determine our habits, our habits determine our character, and our character determines our destiny.",
			"When we put God first, all other things fall into their proper place or drop out of our lives.",
			"If you don’t stand for something, you will fall for anything.",
		},
	}
}

func (r *WordOfWisdomQuoteProvider) Provide() Quote {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(r.phrases))
	return r.phrases[index]
}
