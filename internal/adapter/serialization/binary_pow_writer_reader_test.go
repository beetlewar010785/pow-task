package serialization

import (
	"bytes"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBinaryPOWWriterReader(t *testing.T) {
	testCase := []struct {
		testName string
		pow      message.POW
	}{
		{
			testName: "request POW",
			pow: message.NewPOWRequest(
				"some-challenge",
			),
		},
		{
			testName: "response POW",
			pow: message.NewPOWResponse(
				"some-challenge",
				"some-nonce",
			),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.testName, func(t *testing.T) {
			buf := &bytes.Buffer{}
			sut := NewBinaryPOWWriterReader(buf, buf)

			require.NoError(t, sut.Write(tc.pow))

			readPOW, err := sut.Read()
			require.NoError(t, err)

			assert.Equal(t, tc.pow, readPOW)
		})
	}
}
