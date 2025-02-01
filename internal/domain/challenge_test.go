package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
			nonce:      "776",
			difficulty: 3,
			valid:      true,
		},
		{
			challenge:  "first-challenge",
			nonce:      "123",
			difficulty: 3,
			valid:      false,
		},
		{
			challenge:  "first-challenge",
			nonce:      "776",
			difficulty: 4,
			valid:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf(
			"nonce %s for challenge %s is valid: %v",
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
		const challengeLength = 16
		sut := NewSimpleChallengeRandomizer()

		challenge1 := sut.Generate(challengeLength)
		challenge2 := sut.Generate(challengeLength)

		assert.NotEqual(t, challenge1, challenge2)
	})

	t.Run("generate challenge of the expected length", func(t *testing.T) {
		const challengeLength = 8
		sut := NewSimpleChallengeRandomizer()

		challenge := sut.Generate(challengeLength)

		assert.Len(t, challenge, challengeLength)
	})
}
