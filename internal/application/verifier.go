package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type VerifierFactory interface {
	Create(readWriter domain.ReadWriter) Verifier
}

type Verifier interface {
	Verify() error
}
