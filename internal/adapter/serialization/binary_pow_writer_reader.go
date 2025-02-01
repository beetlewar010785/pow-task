package serialization

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"io"
	"strings"
)

const powFieldsDelimiter = ":"

type BinaryPOWWriterReader struct {
	writer io.Writer
	reader io.Reader
}

func NewBinaryPOWWriterReader(
	writer io.Writer,
	reader io.Reader,
) *BinaryPOWWriterReader {
	return &BinaryPOWWriterReader{
		writer,
		reader,
	}
}

func (r *BinaryPOWWriterReader) Write(grant message.POW) error {
	msg := string(grant.Challenge)
	if grant.Nonce != nil {
		msg += powFieldsDelimiter + grant.Nonce.String()
	}
	msg += "\n" // TODO: improve

	if _, err := r.writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write pow: %w", err)
	}

	return nil
}

func (r *BinaryPOWWriterReader) Read() (message.POW, error) {
	msg, err := bufio.NewReader(r.reader).ReadString('\n')
	if err != nil {
		return message.POW{}, fmt.Errorf("failed to read pow: %w", err)
	}

	msg = strings.TrimSuffix(msg, "\n")
	fields := strings.Split(msg, powFieldsDelimiter)
	if len(fields) == 0 {
		return message.POW{}, errors.New("no pow fields found")
	}

	var nonce *domain.Nonce
	if len(fields) > 1 {
		n := domain.Nonce(fields[1])
		nonce = &n
	}

	return message.NewPOW(domain.Challenge(fields[0]), nonce), nil
}
