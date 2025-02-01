package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomQuoteProvider(t *testing.T) {
	t.Run("provide any quote", func(t *testing.T) {
		quotes := []Quote{
			"quote 1",
			"quote 2",
			"quote 3",
		}

		sut := NewRandomQuoteProvider(quotes)

		actualQuote := sut.Provide()

		assert.Contains(t, quotes, actualQuote)
	})
}
