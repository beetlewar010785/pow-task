package adapter

import (
	"bytes"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStringReadWriter_POWRequest(t *testing.T) {
	t.Run("read same as write", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewStringReadWriter(buf)

		initialPowRequest := domain.NewPOWRequest(
			"some-challenge",
			2,
		)
		require.NoError(t, sut.WritePOWRequest(initialPowRequest))

		actualPOWRequest, err := sut.ReadPOWRequest()
		require.NoError(t, err)

		assert.Equal(t, initialPowRequest, actualPOWRequest)
	})
}

func TestStringReadWriter_POWResponse(t *testing.T) {
	t.Run("read same as write", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewStringReadWriter(buf)

		initialPowResponse := domain.NewPOWResponse(
			"some-challenge",
			2,
		)
		require.NoError(t, sut.WritePOWResponse(initialPowResponse))

		actualPOWResponse, err := sut.ReadPOWResponse()
		require.NoError(t, err)

		assert.Equal(t, initialPowResponse, actualPOWResponse)
	})
}

func TestStringReadWriter_Grant(t *testing.T) {
	t.Run("success grant", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewStringReadWriter(buf)

		initialGrant := domain.SuccessGrant("some-quote")
		require.NoError(t, sut.WriteGrant(initialGrant))

		actualGrant, err := sut.ReadGrant()
		require.NoError(t, err)

		assert.Equal(t, initialGrant, actualGrant)
	})

	t.Run("failure grant", func(t *testing.T) {
		buf := &bytes.Buffer{}
		sut := NewStringReadWriter(buf)

		initialGrant := domain.FailureGrant()
		require.NoError(t, sut.WriteGrant(initialGrant))

		actualGrant, err := sut.ReadGrant()
		require.NoError(t, err)

		assert.Equal(t, initialGrant, actualGrant)
	})
}
