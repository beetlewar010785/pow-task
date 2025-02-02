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
	logger.Info(fmt.Sprintf("starting POW TCP Server at %s", serverAddress))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	powServer := adapter.StartPOWServer(
		serverAddress,
		domain.NewRandomQuoteProvider(defaults.WordOfWisdomQuotes),
		defaults.ChallengeDifficulty,
		defaults.ChallengeLength,
		defaults.VerificationTimeout,
		logger,
	)

	if err := powServer.Listen(); err != nil {
		logger.Error(fmt.Sprintf("failed to setup POW server: %s", err.Error()))
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := powServer.Run(ctx); err != nil {
			logger.Error(fmt.Sprintf("failed to run POW server: %v", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received, stopping server...")

	wg.Wait()
	logger.Info("server stopped")
}
