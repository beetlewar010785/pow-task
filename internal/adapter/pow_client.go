package adapter

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"net"
	"time"
)

func CreatePOWClient(
	serverAddress string,
	solveTimeout time.Duration,
	logger domain.Logger,
) (net.Conn, application.Solver, error) {
	client, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial tcp: %w", err)
	}

	readWriter := NewReadWriterLoggingDecorator(NewStringReadWriter(client), logger)
	challengeVerifier := domain.NewSHA256ChallengeVerifier()
	solver := application.NewPOWSolver(
		domain.NewIncrementalNonceFinder(challengeVerifier),
		readWriter,
		solveTimeout,
	)

	return client, solver, nil
}
