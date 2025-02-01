package application

import (
	"errors"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/beetlewar010785/pow-task/pkg/lib"
)

type POWChallengeHandlerFactory interface {
	Create(in chan []byte, out chan []byte) *POWChallengeHandler
}

type SimplePOWChallengeHandlerFactory struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.GrantProvider
	logger              lib.Logger
}

func NewPOWChallengeHandlerFactory(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.GrantProvider,
	logger lib.Logger,
) *SimplePOWChallengeHandlerFactory {
	return &SimplePOWChallengeHandlerFactory{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		logger,
	}
}

func (r *SimplePOWChallengeHandlerFactory) Create(
	in chan []byte,
	out chan []byte,
) *POWChallengeHandler {
	return NewPOWChallengeHandler(
		in,
		out,
		r.challengeRandomizer,
		r.challengeVerifier,
		r.grantProvider,
		r.logger,
	)
}

type POWChallengeHandler struct {
	in                  chan []byte
	out                 chan []byte
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	granProvider        domain.GrantProvider
	logger              lib.Logger
}

func NewPOWChallengeHandler(
	in chan []byte,
	out chan []byte,
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	granProvider domain.GrantProvider,
	logger lib.Logger,
) *POWChallengeHandler {
	return &POWChallengeHandler{
		in,
		out,
		challengeRandomizer,
		challengeVerifier,
		granProvider,
		logger,
	}
}

func (r *POWChallengeHandler) Handle() error {
	r.logger.Debug("start handling pow challenge")

	challenge := r.challengeRandomizer.Generate()

	r.logger.Debug(fmt.Sprintf("challenge: %s", challenge))

	r.out <- challenge.Bytes()
	nonce := domain.NonceFromBytes(<-r.in)

	r.logger.Debug(fmt.Sprintf("nonce: %s", nonce))

	if !r.challengeVerifier.Verify(challenge, nonce) {
		return errors.New("challenge verification failed")
	}

	grant := r.granProvider.Provide(nonce)
	r.out <- grant.Bytes()

	r.logger.Info(fmt.Sprintf("client granted: %s", grant))

	return nil
}
