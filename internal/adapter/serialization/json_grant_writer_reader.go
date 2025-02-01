package serialization

import (
	"encoding/json"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"io"
)

type JSONGrantWriterReader struct {
	writer io.Writer
	reader io.Reader
}

func NewJSONGrantWriterReader(
	writer io.Writer,
	reader io.Reader,
) *JSONGrantWriterReader {
	return &JSONGrantWriterReader{
		writer,
		reader,
	}
}

func (r *JSONGrantWriterReader) Write(grant message.Grant) error {
	return json.NewEncoder(r.writer).Encode(grant)
}

func (r *JSONGrantWriterReader) Read() (message.Grant, error) {
	var grant message.Grant
	err := json.NewDecoder(r.reader).Decode(&grant)
	return grant, err
}
