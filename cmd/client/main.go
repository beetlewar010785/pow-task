package main

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/caarlos0/env/v10"
	"os"
	"time"
)

type Config struct {
	LogLevel        adapter.LogLevel `env:"LOG_LEVEL" envDefault:"2"`
	SolveTimeoutSec int              `env:"SOLVE_TIMEOUT_SEC" envDefault:"10"`
	ServerAddress   string           `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	logger := adapter.NewStdLogger("client", cfg.LogLevel)

	logger.Info(fmt.Sprintf("env:  %+v", cfg))
	logger.Info(fmt.Sprintf("connecting to %s", cfg.ServerAddress))

	conn, grantReceiver, err := adapter.CreatePOWClient(
		cfg.ServerAddress,
		time.Duration(cfg.SolveTimeoutSec)*time.Second,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to %s", cfg.ServerAddress))
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
