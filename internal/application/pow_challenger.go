package application

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POWChallengerFactory struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.QuoteProvider
	challengeDifficulty domain.Difficulty
}

func NewPOWChallengerFactory(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.QuoteProvider,
	challengeDifficulty domain.Difficulty,
) *POWChallengerFactory {
	return &POWChallengerFactory{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
	}
}

func (r *POWChallengerFactory) Create(readWriter domain.ReadWriter) Challenger {
	return NewPOWChallenger(
		r.challengeRandomizer,
		r.challengeVerifier,
		r.grantProvider,
		r.challengeDifficulty,
		readWriter,
	)
}

type POWChallenger struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.QuoteProvider
	challengeDifficulty domain.Difficulty
	readWriter          domain.ReadWriter
}

func NewPOWChallenger(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.QuoteProvider,
	challengeDifficulty domain.Difficulty,
	readWriter domain.ReadWriter,
) *POWChallenger {
	return &POWChallenger{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
		readWriter,
	}
}

func (r *POWChallenger) Challenge() error {
	challenge := r.challengeRandomizer.Generate()

	powRequest := domain.NewPOWRequest(challenge, r.challengeDifficulty)
	if err := r.readWriter.WritePowRequest(powRequest); err != nil {
		return fmt.Errorf("error writing pow request: %w", err)
	}

	powResponse, err := r.readWriter.ReadPowResponse()
	if err != nil {
		return fmt.Errorf("error reading pow response: %w", err)
	}

	var grant domain.Grant
	if !r.challengeVerifier.Verify(challenge, powResponse.Nonce, r.challengeDifficulty) {
		grant = domain.FailureGrant()
	} else {
		quote := r.grantProvider.Provide()
		grant = domain.SuccessGrant(string(quote))
	}

	if err := r.readWriter.WriteGrant(grant); err != nil {
		return fmt.Errorf("error writing grant: %w", err)
	}

	return nil
}
