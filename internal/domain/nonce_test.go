package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementalNonceFinder(t *testing.T) {
	testCases := []struct {
		challenge     Challenge
		expectedNonce Nonce
		difficulty    Difficulty
	}{
		{
			"first-challenge",
			"776",
			3,
		},
		{
			"second-challenge",
			"14707",
			4,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s nonce for challenge %s", tc.expectedNonce, tc.challenge), func(t *testing.T) {
			sut := NewIncrementalNonceFinder()

			actualNonce := sut.Find(tc.challenge, tc.difficulty)

			assert.Equal(t, tc.expectedNonce, actualNonce)
		})
	}
}
