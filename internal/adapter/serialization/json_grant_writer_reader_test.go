package serialization

import (
	"bytes"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONGrantWriterReader(t *testing.T) {
	t.Run("success grant", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewJSONGrantWriterReader(buf, buf)

		initialGrant := message.SuccessGrant("some-quote")
		require.NoError(t, sut.Write(initialGrant))

		readGrant, err := sut.Read()
		require.NoError(t, err)

		assert.Equal(t, initialGrant, readGrant)
	})

	t.Run("failure grant", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewJSONGrantWriterReader(buf, buf)

		initialGrant := message.FailureGrant()
		require.NoError(t, sut.Write(initialGrant))

		readGrant, err := sut.Read()
		require.NoError(t, err)

		assert.Equal(t, initialGrant, readGrant)
	})
}
