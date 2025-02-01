package application

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter/serialization"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"io"
)

type POWChallengerFactory struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.GrantProvider
	challengeDifficulty domain.Difficulty
	challengeLength     int
}

func NewPOWChallengerFactory(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.GrantProvider,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
) *POWChallengerFactory {
	return &POWChallengerFactory{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
		challengeLength,
	}
}

func (r *POWChallengerFactory) Create(
	reader io.Reader,
	writer io.Writer,
) Challenger {
	powWriterReader := serialization.NewBinaryPOWWriterReader(
		writer,
		reader,
	)

	grantWriterReader := serialization.NewJSONGrantWriterReader(
		writer,
		reader,
	)

	return NewPOWChallenger(
		r.challengeRandomizer,
		r.challengeVerifier,
		r.grantProvider,
		r.challengeDifficulty,
		r.challengeLength,
		powWriterReader,
		powWriterReader,
		grantWriterReader,
	)
}

type POWChallenger struct {
	challengeRandomizer domain.ChallengeRandomizer
	challengeVerifier   domain.ChallengeVerifier
	grantProvider       domain.GrantProvider
	challengeDifficulty domain.Difficulty
	challengeLength     int
	powWriter           message.POWWriter
	powReader           message.POWReader
	grantWriter         message.GrantWriter
}

func NewPOWChallenger(
	challengeRandomizer domain.ChallengeRandomizer,
	challengeVerifier domain.ChallengeVerifier,
	grantProvider domain.GrantProvider,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
	powWriter message.POWWriter,
	powReader message.POWReader,
	grantWriter message.GrantWriter,
) *POWChallenger {
	return &POWChallenger{
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
		challengeLength,
		powWriter,
		powReader,
		grantWriter,
	}
}

func (r *POWChallenger) Challenge() error {
	challenge := r.challengeRandomizer.Generate(r.challengeLength)

	if err := r.powWriter.Write(message.NewPOWRequest(challenge)); err != nil {
		return fmt.Errorf("error writing pow request: %w", err)
	}

	pow, err := r.powReader.Read()
	if err != nil {
		return fmt.Errorf("error reading pow response: %w", err)
	}

	if pow.Nonce == nil {
		return fmt.Errorf("nonce is missing")
	}

	var grant message.Grant
	if !r.challengeVerifier.Verify(challenge, *pow.Nonce, r.challengeDifficulty) {
		grant = message.FailureGrant()
	} else {
		quote := r.grantProvider.Provide()
		grant = message.SuccessGrant(string(quote))
	}

	if err := r.grantWriter.Write(grant); err != nil {
		return fmt.Errorf("error writing grant: %w", err)
	}

	return nil
}
