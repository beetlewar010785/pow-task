package message

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POW struct {
	Challenge domain.Challenge
	Nonce     *domain.Nonce // nil for request
}

func NewPOW(challenge domain.Challenge, nonce *domain.Nonce) POW {
	return POW{
		challenge,
		nonce,
	}
}

func NewPOWRequest(challenge domain.Challenge) POW {
	return NewPOW(challenge, nil)
}

func NewPOWResponse(challenge domain.Challenge, nonce domain.Nonce) POW {
	return NewPOW(challenge, &nonce)
}

func (r POW) WithNonce(nonce domain.Nonce) POW {
	return NewPOWResponse(r.Challenge, nonce)
}

type POWWriter interface {
	Write(pow POW) error
}

type POWReader interface {
	Read() (POW, error)
}
