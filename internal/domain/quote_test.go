package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordOfWisdomQuoteProvider(t *testing.T) {
	t.Run("provide any quote", func(t *testing.T) {
		sut := NewWordOfWisdomQuoteProvider()

		actualQuote := sut.Provide()

		assert.NotEmpty(t, actualQuote, actualQuote)
	})
}
