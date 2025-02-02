package main

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/defaults"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	serverAddress := fmt.Sprintf(":%s", defaults.ServerPort)
	logger := adapter.NewStdLogger("server", defaults.LogLevel)
	logger.Info(fmt.Sprintf("starting TCP Server at %s", serverAddress))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tcpServer := adapter.StartTCPServer(
		serverAddress,
		domain.NewRandomQuoteProvider(defaults.WordOfWisdomQuotes),
		defaults.ChallengeDifficulty,
		defaults.ChallengeLength,
		defaults.VerificationTimeout,
		logger,
	)

	if err := tcpServer.Listen(); err != nil {
		logger.Error(fmt.Sprintf("failed to setup tcp server: %s", err.Error()))
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tcpServer.Run(ctx); err != nil {
			logger.Error(fmt.Sprintf("failed to run  tcp server: %v", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received, stopping server...")

	wg.Wait()
	logger.Info("server stopped")
}
