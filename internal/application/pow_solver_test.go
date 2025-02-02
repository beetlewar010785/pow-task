package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPOWSolver(t *testing.T) {
	t.Run("solve POW and get grant", func(t *testing.T) {
		nonce := domain.Nonce(123)
		challenge := domain.Challenge("some-challenge")
		difficulty := domain.Difficulty(234)
		grant := domain.SuccessGrant("some-quote")

		nonceFinder := new(nonceFinderMock)
		nonceFinder.withNonce(nonce, challenge, difficulty)

		readWriter := new(domain.ReadWriterMock)
		readWriter.POWRequest = domain.NewPOWRequest(challenge, difficulty)
		readWriter.Grant = grant

		sut := NewPOWSolver(nonceFinder, readWriter)

		actualGrant, err := sut.Solve()
		require.NoError(t, err)

		assert.Equal(t, domain.NewPOWResponse(challenge, nonce), readWriter.POWResponse)
		assert.Equal(t, grant, actualGrant)
	})
}

type nonceFinderMock struct {
	expectedChallenge  domain.Challenge
	expectedDifficulty domain.Difficulty
	nonceToFind        domain.Nonce
}

func (r *nonceFinderMock) withNonce(
	actualNonce domain.Nonce,
	expectedChallenge domain.Challenge,
	expectedDifficulty domain.Difficulty,
) {
	r.expectedChallenge = expectedChallenge
	r.expectedDifficulty = expectedDifficulty
	r.nonceToFind = actualNonce
}

func (r *nonceFinderMock) Find(
	challenge domain.Challenge,
	difficulty domain.Difficulty,
) domain.Nonce {
	if r.expectedChallenge == challenge && r.expectedDifficulty == difficulty {
		return r.nonceToFind
	}

	return 0
}
