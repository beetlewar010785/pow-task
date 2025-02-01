package main

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/beetlewar010785/pow-task/pkg/lib"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	const serverAddress = ":8080"
	const challengeDifficulty = 4
	const challengeLength = 16
	wordOfWisdomQuotes := []domain.Grant{
		"Cease to be idle; cease to be unclean; cease to find fault one with another.",
		"A man is saved no faster than he gains knowledge.",
		"Our thoughts determine our actions, our actions determine our habits, our habits determine our character, and our character determines our destiny.",
		"When we put God first, all other things fall into their proper place or drop out of our lives.",
		"If you don’t stand for something, you will fall for anything.",
	}

	logger := lib.NewStdLogger("server", lib.LogLevelInfo)
	logger.Info(fmt.Sprintf("starting TCP Server at %s", serverAddress))

	ctx, stop := createSignalContext()
	defer stop()

	tcpServer, err := setupServer(
		serverAddress,
		challengeDifficulty,
		challengeLength,
		wordOfWisdomQuotes,
		logger,
	)
	if err != nil {
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

func createSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func setupServer(
	serverAddress string,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
	wordOfWisdomQuotes []domain.Grant,
	logger lib.Logger,
) (*adapter.TCPServer, error) {
	challengeRandomizer := domain.NewSimpleChallengeRandomizer()
	challengeVerifier := domain.NewSimpleChallengeVerifier()
	grantProvider := domain.NewRandomPhraseGrantProvider(wordOfWisdomQuotes)
	powServerFactory := application.NewPOWChallengeHandlerFactory(
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
		challengeLength,
		logger,
	)

	tcpServer := adapter.NewTCPServer(serverAddress, powServerFactory, logger)
	if err := tcpServer.Listen(); err != nil {
		return nil, fmt.Errorf("failed to listen tcp server: %w", err)
	}

	return tcpServer, nil
}
