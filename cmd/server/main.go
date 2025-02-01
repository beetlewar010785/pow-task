package main

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/domain"
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

	logger := adapter.NewStdLogger("server", adapter.LogLevelInfo)
	logger.Info(fmt.Sprintf("starting TCP Server at %s", serverAddress))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	tcpServer := adapter.StartTCPServer(
		serverAddress,
		domain.NewRandomPhraseGrantProvider(wordOfWisdomQuotes),
		challengeDifficulty,
		challengeLength,
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
