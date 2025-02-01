package domain

import (
	"math/rand"
	"time"
)

type Grant string

func (r Grant) Bytes() []byte {
	return []byte(r)
}

type GrantProvider interface {
	Provide() Grant
}

type RandomPhraseGrantProvider struct {
	phrases []Grant
}

func NewRandomPhraseGrantProvider(grants []Grant) *RandomPhraseGrantProvider {
	return &RandomPhraseGrantProvider{
		grants,
	}
}

func (r *RandomPhraseGrantProvider) Provide() Grant {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rnd.Intn(len(r.phrases))
	return r.phrases[index]
}
