package domain

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleChallengeVerifier(t *testing.T) {
	testCases := []struct {
		challenge  Challenge
		nonce      Nonce
		difficulty Difficulty
		valid      bool
	}{
		{
			challenge:  "first-challenge",
			nonce:      776,
			difficulty: 3,
			valid:      true,
		},
		{
			challenge:  "first-challenge",
			nonce:      123,
			difficulty: 3,
			valid:      false,
		},
		{
			challenge:  "first-challenge",
			nonce:      776,
			difficulty: 4,
			valid:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf(
			"nonce %d for challenge %s is valid: %v",
			tc.nonce,
			tc.challenge,
			tc.valid,
		), func(t *testing.T) {
			sut := NewSimpleChallengeVerifier()

			isActualValid := sut.Verify(tc.challenge, tc.nonce, tc.difficulty)

			assert.Equal(t, tc.valid, isActualValid)
		})
	}
}

func TestSimpleChallengeRandomizer(t *testing.T) {
	t.Run("generate random challenge", func(t *testing.T) {
		sut := NewSimpleChallengeRandomizer(12)

		challenge1 := sut.Generate()
		challenge2 := sut.Generate()

		assert.NotEqual(t, challenge1, challenge2)
	})

	t.Run("generate challenge of the expected length", func(t *testing.T) {
		sut := NewSimpleChallengeRandomizer(4)

		challenge := sut.Generate()

		assert.Len(t, challenge, 4)
	})
}

func TestIncrementalNonceFinder(t *testing.T) {
	t.Run("incrementally search for nonce", func(t *testing.T) {
		challengeVerifier := new(challengeVerifierMock)
		challengeVerifier.expectedNonce = 42
		sut := NewIncrementalNonceFinder(challengeVerifier)

		actualNonce, err := sut.Find(context.Background(), "any", 3)
		require.NoError(t, err)

		assert.Equal(t, challengeVerifier.expectedNonce, actualNonce)
	})
}

type challengeVerifierMock struct {
	expectedNonce Nonce
}

func (r *challengeVerifierMock) Verify(
	_ Challenge,
	nonce Nonce,
	_ Difficulty,
) bool {
	return nonce == r.expectedNonce
}
