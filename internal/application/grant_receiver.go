package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type GrantReceiver interface {
	Receive() (domain.Grant, error)
}
