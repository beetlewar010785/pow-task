package main

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/defaults"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"os"
)

func main() {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = fmt.Sprintf("localhost:%s", defaults.ServerPort)
	}

	logger := adapter.NewStdLogger("client", defaults.LogLevel)
	logger.Info(fmt.Sprintf("connecting to %s", serverAddress))

	conn, grantReceiver, err := adapter.CreateTCPClient(serverAddress, defaults.VerificationTimeout)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to %s", serverAddress))
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	grant, err := grantReceiver.Solve()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to receive grant message: %s", err))
		os.Exit(1)
	}

	if grant.Result != domain.GrantResultSuccess {
		logger.Error(fmt.Sprintf("grant failed: %d", grant.Result))
		os.Exit(1)
	}

	logger.Info(string(*grant.Quote))
}
