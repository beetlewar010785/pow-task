package application

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POWGrantReceiver struct {
	nonceFinder domain.NonceFinder
	difficulty  domain.Difficulty
	readWriter  domain.ReadWriter
}

func NewPOWGrantReceiver(
	nonceFinder domain.NonceFinder,
	difficulty domain.Difficulty,
	readWriter domain.ReadWriter,
) *POWGrantReceiver {
	return &POWGrantReceiver{
		nonceFinder,
		difficulty,
		readWriter,
	}
}

func (r *POWGrantReceiver) Receive() (domain.Grant, error) {
	powRequest, err := r.readWriter.ReadPowRequest()
	if err != nil {
		return domain.Grant{}, fmt.Errorf("error reading pow request: %w", err)
	}

	nonce := r.nonceFinder.Find(powRequest.Challenge, r.difficulty)
	powResponse := domain.NewPOWResponse(powRequest.Challenge, nonce)

	if err := r.readWriter.WritePowResponse(powResponse); err != nil {
		return domain.Grant{}, fmt.Errorf("error writing pow response: %w", err)
	}

	grant, err := r.readWriter.ReadGrant()
	if err != nil {
		return domain.Grant{}, fmt.Errorf("error reading grant: %w", err)
	}

	return grant, nil
}
