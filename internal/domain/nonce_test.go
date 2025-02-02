package domain

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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
