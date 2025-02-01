package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomQuoteProvider(t *testing.T) {
	t.Run("provide any quote", func(t *testing.T) {
		phrases := []Quote{
			"quote 1",
			"quote 2",
			"quote 3",
		}

		sut := NewRandomPhraseGrantProvider(phrases)

		actualGrant := sut.Provide()

		assert.Contains(t, phrases, actualGrant)
	})
}
