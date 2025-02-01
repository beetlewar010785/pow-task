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
	"syscall"
)

func main() {
	const serverAddress = "localhost:8080"
	const difficulty = 4

	logger := lib.NewStdLogger("client", lib.LogLevelInfo)
	logger.Info(fmt.Sprintf("connecting to %s", serverAddress))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	in := make(chan []byte)
	out := make(chan []byte)

	tcpClient := adapter.NewTCPClient(serverAddress, in, out, logger)
	go func() {
		if err := tcpClient.Run(ctx); err != nil {
			logger.Error(fmt.Sprintf("tcp client run error: %v", err))
			os.Exit(1)
		}
	}()

	nonceFinder := domain.NewIncrementalNonceFinder()
	powClient := application.NewPOWGrantReceiver(in, out, nonceFinder, difficulty, logger)

	resultChan := make(chan domain.Grant)
	go func() {
		resultChan <- powClient.Receive()
	}()

	select {
	case <-ctx.Done():
		logger.Warn("client shutting down due to interrupt signal")
	case grant := <-resultChan:
		logger.Info(fmt.Sprintf("received grant: %s", grant))
	}
}
