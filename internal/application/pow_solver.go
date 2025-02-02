package application

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"time"
)

type POWSolver struct {
	nonceFinder      domain.NonceFinder
	readWriter       domain.ReadWriter
	findNonceTimeout time.Duration
}

func NewPOWSolver(
	nonceFinder domain.NonceFinder,
	readWriter domain.ReadWriter,
	findNonceTimeout time.Duration,
) *POWSolver {
	return &POWSolver{
		nonceFinder,
		readWriter,
		findNonceTimeout,
	}
}

func (r *POWSolver) Solve() (domain.Grant, error) {
	powRequest, err := r.readWriter.ReadPOWRequest()
	if err != nil {
		return domain.Grant{}, fmt.Errorf("error reading pow request: %w", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), r.findNonceTimeout)
	defer cancel()

	nonce, err := r.nonceFinder.Find(ctxWithTimeout, powRequest.Challenge, powRequest.Difficulty)
	if err != nil {
		return domain.Grant{}, fmt.Errorf("error finding nonce: %w", err)
	}

	powResponse := domain.NewPOWResponse(powRequest.Challenge, nonce)

	if err := r.readWriter.WritePOWResponse(powResponse); err != nil {
		return domain.Grant{}, fmt.Errorf("error writing pow response: %w", err)
	}

	grant, err := r.readWriter.ReadGrant()
	if err != nil {
		return domain.Grant{}, fmt.Errorf("error reading grant: %w", err)
	}

	return grant, nil
}
