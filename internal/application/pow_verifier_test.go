package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPOWVerifier(t *testing.T) {
	t.Run("success grant when verified POW", func(t *testing.T) {
		nonce := domain.Nonce(123)
		challenge := domain.Challenge("some-challenge")
		difficulty := domain.Difficulty(234)
		quote := domain.Quote("some-quote")

		challengeRandomizer := new(challengeRandomizerMock)
		challengeRandomizer.challenge = challenge

		challengeVerifier := new(challengeVerifierMock)
		challengeVerifier.verified(challenge, nonce, difficulty)

		quoteProvider := new(quoteProviderMock)
		quoteProvider.quote = quote

		readWriter := new(domain.ReadWriterMock)
		readWriter.POWResponse = domain.NewPOWResponse(challenge, nonce)

		sut := NewPOWVerifier(
			challengeRandomizer,
			challengeVerifier,
			quoteProvider,
			difficulty,
			readWriter,
		)

		err := sut.Verify()
		require.NoError(t, err)

		assert.Equal(t, domain.NewPOWRequest(challenge, difficulty), readWriter.POWRequest)
		assert.Equal(t, domain.SuccessGrant(quote), readWriter.Grant)
	})

	t.Run("failure grant when non verified POW", func(t *testing.T) {
		challenge := domain.Challenge("some-challenge")

		challengeRandomizer := new(challengeRandomizerMock)
		challengeRandomizer.challenge = challenge

		readWriter := new(domain.ReadWriterMock)
		readWriter.POWResponse = domain.NewPOWResponse(challenge, domain.Nonce(123))

		sut := NewPOWVerifier(
			challengeRandomizer,
			new(challengeVerifierMock),
			new(quoteProviderMock),
			domain.Difficulty(234),
			readWriter,
		)

		err := sut.Verify()
		require.NoError(t, err)

		assert.Equal(t, domain.FailureGrant(), readWriter.Grant)
	})
}

type challengeRandomizerMock struct {
	challenge domain.Challenge
}

func (r *challengeRandomizerMock) Generate() domain.Challenge {
	return r.challenge
}

type challengeVerifierMock struct {
	expectedChallenge  domain.Challenge
	expectedNonce      domain.Nonce
	expectedDifficulty domain.Difficulty
}

func (r *challengeVerifierMock) verified(
	challenge domain.Challenge,
	nonce domain.Nonce,
	difficulty domain.Difficulty,
) {
	r.expectedChallenge = challenge
	r.expectedNonce = nonce
	r.expectedDifficulty = difficulty
}

func (r *challengeVerifierMock) Verify(
	challenge domain.Challenge,
	nonce domain.Nonce,
	difficulty domain.Difficulty,
) bool {
	return challenge == r.expectedChallenge &&
		nonce == r.expectedNonce &&
		difficulty == r.expectedDifficulty
}

type quoteProviderMock struct {
	quote domain.Quote
}

func (r *quoteProviderMock) Provide() domain.Quote {
	return r.quote
}
