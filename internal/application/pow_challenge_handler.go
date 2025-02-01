package application

import (
	"errors"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POWChallengeHandlerFactory interface {
	Create(in chan []byte, out chan []byte) *POWChallengeHandler
}

type SimplePOWChallengeHandlerFactory struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.GrantProvider
	challengeDifficulty domain.Difficulty
	challengeLength     int
	logger              domain.Logger
}

func NewPOWChallengeHandlerFactory(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.GrantProvider,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
	logger domain.Logger,
) *SimplePOWChallengeHandlerFactory {
	return &SimplePOWChallengeHandlerFactory{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
		challengeLength,
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
		r.challengeDifficulty,
		r.challengeLength,
		r.logger,
	)
}

type POWChallengeHandler struct {
	in                  chan []byte
	out                 chan []byte
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	granProvider        domain.GrantProvider
	challengeDifficulty domain.Difficulty
	challengeLength     int
	logger              domain.Logger
}

func NewPOWChallengeHandler(
	in chan []byte,
	out chan []byte,
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	granProvider domain.GrantProvider,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
	logger domain.Logger,
) *POWChallengeHandler {
	return &POWChallengeHandler{
		in,
		out,
		challengeRandomizer,
		challengeVerifier,
		granProvider,
		challengeDifficulty,
		challengeLength,
		logger,
	}
}

func (r *POWChallengeHandler) Handle() error {
	r.logger.Debug("start handling pow challenge")

	challenge := r.challengeRandomizer.Generate(r.challengeLength)

	r.logger.Debug(fmt.Sprintf("challenge: %s", challenge))

	r.out <- challenge.Bytes()
	nonce := domain.NonceFromBytes(<-r.in)

	r.logger.Debug(fmt.Sprintf("nonce: %s", nonce))

	if !r.challengeVerifier.Verify(challenge, nonce, r.challengeDifficulty) {
		return errors.New("challenge verification failed")
	}

	grant := r.granProvider.Provide()
	r.out <- grant.Bytes()

	r.logger.Info(fmt.Sprintf("client granted: %s", grant))

	return nil
}
