package test

import (
	"context"
	"testing"
	"time"

	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/beetlewar010785/pow-task/pkg/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Run("should pass PoW challenge", func(t *testing.T) {
		const randomPort = ":0"
		const difficulty = 5
		logLevel := lib.LogLevelDebug

		grantProvider := &grantProviderMock{grant: "expected-grant"}
		tcpServer := createTCPServer(randomPort, grantProvider, lib.NewStdLogger("server", logLevel))

		require.NoError(t, tcpServer.Listen())

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		go func() {
			err := tcpServer.Run(ctx)
			assert.NoError(t, err, "server.Run() returned an unexpected error")
		}()

		WaitForServer(t, tcpServer)

		serverAddress := tcpServer.Address()
		powClient, tcpClient, err := createClient(
			serverAddress,
			difficulty,
			lib.NewStdLogger("client", logLevel),
		)
		require.NoError(t, err)

		go func() {
			err := tcpClient.Run(ctx)
			assert.NoError(t, err, "server.Run() returned an unexpected error")
		}()

		actualGrant := powClient.Receive()
		assert.Equal(t, grantProvider.grant, actualGrant)
	})
}

func WaitForServer(t *testing.T, server *adapter.TCPServer) {
	timeout := time.After(2 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Server did not start listening in time")
		case <-ticker.C:
			if server.Address() != "" {
				return
			}
		}
	}
}

func createTCPServer(
	serverAddress string,
	grantProvider domain.GrantProvider,
	logger lib.Logger,
) *adapter.TCPServer {
	challengeRandomizer := domain.NewSimpleChallengeRandomizer()
	challengeVerifier := domain.NewSimpleChallengeVerifier()
	powServerFactory := application.NewPOWChallengeHandlerFactory(challengeRandomizer, challengeVerifier, grantProvider, logger)
	return adapter.NewTCPServer(serverAddress, powServerFactory, logger)
}

func createClient(
	serverAddress string,
	difficulty domain.Difficulty,
	logger lib.Logger,
) (*application.POWGrantReceiver, *adapter.TCPClient, error) {
	in := make(chan []byte)
	out := make(chan []byte)

	tcpClient := adapter.NewTCPClient(
		serverAddress,
		in,
		out,
		logger,
	)

	return application.NewPOWGrantReceiver(
		in,
		out,
		domain.NewIncrementalNonceFinder(),
		difficulty,
		logger,
	), tcpClient, nil
}

type grantProviderMock struct {
	grant domain.Grant
}

func (r *grantProviderMock) Provide(_ domain.Nonce) domain.Grant {
	return r.grant
}
