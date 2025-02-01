package adapter

import (
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/adapter/serialization"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"net"
)

func CreateTCPClient(
	serverAddress string,
	challengeDifficulty domain.Difficulty,
) (net.Conn, application.GrantReceiver, error) {
	client, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial tcp: %w", err)
	}

	powWriterReader := serialization.NewBinaryPOWWriterReader(client, client)
	grantReceiver := application.NewPOWGrantReceiver(
		domain.NewIncrementalNonceFinder(),
		challengeDifficulty,
		powWriterReader,
		powWriterReader,
		serialization.NewJSONGrantWriterReader(client, client),
	)

	return client, grantReceiver, nil
}
