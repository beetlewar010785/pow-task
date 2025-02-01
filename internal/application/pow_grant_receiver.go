package application

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type POWGrantReceiver struct {
	in          chan []byte
	out         chan []byte
	nonceFinder domain.NonceFinder
	difficulty  domain.Difficulty
	logger      domain.Logger
}

func NewPOWGrantReceiver(
	in chan []byte,
	out chan []byte,
	nonceFinder domain.NonceFinder,
	difficulty domain.Difficulty,
	logger domain.Logger,
) *POWGrantReceiver {
	return &POWGrantReceiver{
		in,
		out,
		nonceFinder,
		difficulty,
		logger,
	}
}

func (r *POWGrantReceiver) Receive() domain.Grant {
	r.logger.Debug("start handling pow client")

	challenge := domain.ChallengeFromBytes(<-r.in)

	r.logger.Debug(fmt.Sprintf("received challenge: %s", challenge))

	nonce := r.nonceFinder.Find(challenge, r.difficulty)

	r.logger.Debug(fmt.Sprintf("found nonce: %s", nonce))

	r.out <- nonce.Bytes()

	grant := domain.Grant(<-r.in)

	r.logger.Debug(fmt.Sprintf("received grant: %s", grant))

	return grant
}
