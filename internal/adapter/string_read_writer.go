package adapter

import (
	"bufio"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"io"
	"strconv"
	"strings"
)

const powFieldDelimiter string = ":"

type StringReadWriter struct {
	readWriter io.ReadWriter
}

func NewStringReadWriter(readWriter io.ReadWriter) *StringReadWriter {
	return &StringReadWriter{
		readWriter,
	}
}

func (r *StringReadWriter) WritePowRequest(powRequest domain.POWRequest) error {
	return r.write(string(powRequest.Challenge), strconv.Itoa(int(powRequest.Difficulty)))
}

func (r *StringReadWriter) ReadPowRequest() (domain.POWRequest, error) {
	fields, err := r.read()
	if err != nil {
		return domain.POWRequest{}, err
	}

	i, _ := strconv.Atoi(fields[1])
	return domain.NewPOWRequest(domain.Challenge(fields[0]), domain.Difficulty(i)), nil
}

func (r *StringReadWriter) WritePowResponse(powResponse domain.POWResponse) error {
	return r.write(string(powResponse.Challenge), strconv.Itoa(int(powResponse.Nonce)))
}

func (r *StringReadWriter) ReadPowResponse() (domain.POWResponse, error) {
	fields, err := r.read()
	if err != nil {
		return domain.POWResponse{}, err
	}

	i, _ := strconv.Atoi(fields[1])
	return domain.NewPOWResponse(domain.Challenge(fields[0]), domain.Nonce(i)), nil
}

func (r *StringReadWriter) WriteGrant(grant domain.Grant) error {
	fields := []string{
		strconv.Itoa(int(grant.Result)),
	}

	if grant.Quote != nil {
		fields = append(fields, *grant.Quote)
	}

	return r.write(fields...)
}

func (r *StringReadWriter) ReadGrant() (domain.Grant, error) {
	fields, err := r.read()
	if err != nil {
		return domain.Grant{}, err
	}

	var quote *string
	if len(fields) > 1 {
		quote = &fields[1]
	}

	i, _ := strconv.Atoi(fields[0])
	return domain.NewGrant(domain.GrantResult(i), quote), nil
}

func (r *StringReadWriter) write(params ...string) error {
	s := ""
	for _, param := range params {
		if len(s) > 0 {
			s += powFieldDelimiter
		}
		s += param
	}

	_, err := fmt.Fprintln(r.readWriter, s)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

func (r *StringReadWriter) read() ([]string, error) {
	rw := bufio.NewReader(r.readWriter)
	s, err := rw.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}
	s = strings.TrimSuffix(s, "\n")

	return strings.Split(s, powFieldDelimiter), nil
}
