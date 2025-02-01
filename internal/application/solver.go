package application

import (
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type Solver interface {
	Solve() (domain.Grant, error)
}
