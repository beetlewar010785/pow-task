package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomPhraseGrantProvider(t *testing.T) {
	t.Run("provide any grant from the phrases", func(t *testing.T) {
		phrases := []Grant{
			"phrase 1",
			"phrase 2",
			"phrase 3",
		}

		sut := NewRandomPhraseGrantProvider(phrases)

		actualGrant := sut.Provide()

		assert.Contains(t, phrases, actualGrant)
	})
}
