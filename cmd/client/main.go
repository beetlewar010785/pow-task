package main

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/beetlewar010785/pow-task/pkg/lib"
	"os"
)

func main() {
	logger := lib.NewStdLogger("client")

	const serverAddress = "localhost:8080"
	const difficulty = 4

	incomingChannel := make(chan []byte)
	outgoingChannel := make(chan []byte)

	connection, err := adapter.CreateAndRunTCPConnection(serverAddress, incomingChannel, outgoingChannel, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(connection *adapter.TCPConnection) {
		_ = connection.Close()
	}(connection)

	nonceFinder := domain.NewIncrementalNonceFinder()

	powClient := application.NewPoWClient(
		incomingChannel,
		outgoingChannel,
		nonceFinder,
		difficulty,
		logger,
	)

	grant := powClient.Handle()
	logger.Info(fmt.Sprintf("received grant: %s", grant))
}
