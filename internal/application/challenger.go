package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type ChallengerFactory interface {
	Create(readWriter domain.ReadWriter) Challenger
}

type Challenger interface {
	Challenge() error
}
