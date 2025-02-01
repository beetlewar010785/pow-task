package application

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POWGrantReceiver struct {
	nonceFinder domain.NonceFinder
	difficulty  domain.Difficulty
	powWriter   message.POWWriter
	powReader   message.POWReader
	grantReader message.GrantReader
}

func NewPOWGrantReceiver(
	nonceFinder domain.NonceFinder,
	difficulty domain.Difficulty,
	powWriter message.POWWriter,
	powReader message.POWReader,
	grantReader message.GrantReader,
) *POWGrantReceiver {
	return &POWGrantReceiver{
		nonceFinder,
		difficulty,
		powWriter,
		powReader,
		grantReader,
	}
}

func (r *POWGrantReceiver) Receive() (message.Grant, error) {
	pow, err := r.powReader.Read()
	if err != nil {
		return message.Grant{}, fmt.Errorf("failed to read pow: %w", err)
	}

	nonce := r.nonceFinder.Find(pow.Challenge, r.difficulty)
	powWithNonce := pow.WithNonce(nonce)

	if err := r.powWriter.Write(powWithNonce); err != nil {
		return message.Grant{}, fmt.Errorf("failed to write pow: %w", err)
	}

	grant, err := r.grantReader.Read()
	if err != nil {
		return message.Grant{}, fmt.Errorf("failed to read grant: %w", err)
	}

	return grant, nil
}
