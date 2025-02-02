package main

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/caarlos0/env/v10"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	Port                   int               `env:"PORT" envDefault:"8080"`
	LogLevel               adapter.LogLevel  `env:"LOG_LEVEL" envDefault:"2"`
	ChallengeDifficulty    domain.Difficulty `env:"CHALLENGE_DIFFICULTY" envDefault:"4"`
	ChallengeLength        int               `env:"CHALLENGE_LENGTH" envDefault:"16"`
	VerificationTimeoutSec int               `env:"VERIFICATION_TIMEOUT_SEC" envDefault:"10"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	logger := adapter.NewStdLogger("server", cfg.LogLevel)
	logger.Info(fmt.Sprintf("env:  %+v", cfg))

	logger.Info(fmt.Sprintf("starting POW TCP Server at %d", cfg.Port))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverAddress := fmt.Sprintf(":%d", cfg.Port)
	powServer := adapter.StartPOWServer(
		serverAddress,
		domain.NewWordOfWisdomQuoteProvider(),
		cfg.ChallengeDifficulty,
		cfg.ChallengeLength,
		time.Duration(cfg.VerificationTimeoutSec)*time.Second,
		logger,
	)

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
