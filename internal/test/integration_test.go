package test

import (
	"context"
	"github.com/beetlewar010785/pow-task/internal/application/message"
	"net"
	"testing"
	"time"

	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	type testSuite struct {
		grantProviderMock *grantProviderMock
		server            *adapter.TCPServer
		client            net.Conn
		grantReceiver     application.GrantReceiver
		ctx               context.Context
	}

	setup := func(
		t *testing.T,
		challengeDifficulty domain.Difficulty,
		challengeLength int,
	) testSuite {
		const randomPort = ":0"

		grantProvider := new(grantProviderMock)
		tcpServer := adapter.StartTCPServer(
			randomPort,
			grantProvider,
			challengeDifficulty,
			challengeLength,
			adapter.NewStdLogger("server", adapter.LogLevelInfo),
		)

		require.NoError(t, tcpServer.Listen())

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		go func() {
			err := tcpServer.Run(ctx)
			assert.NoError(t, err, "server.Run() returned an unexpected error")
		}()

		WaitForServer(t, tcpServer)

		conn, grantReceiver, err := adapter.CreateTCPClient(tcpServer.Address(), challengeDifficulty)
		require.NoError(t, err)

		return testSuite{
			grantProvider,
			tcpServer,
			conn,
			grantReceiver,
			ctx,
		}
	}

	tearDown := func(suite testSuite) {
		suite.ctx.Done()
		_ = suite.client.Close()
	}

	t.Run("should pass PoW challenge", func(t *testing.T) {
		const challengeDifficulty = 4
		const challengeLength = 16

		suite := setup(t, challengeDifficulty, challengeLength)
		defer tearDown(suite)

		suite.grantProviderMock.grant = "expected-grant"
		actualGrant, err := suite.grantReceiver.Receive()
		require.NoError(t, err)

		expectedGrant := message.SuccessGrant(string(suite.grantProviderMock.grant))
		assert.Equal(t, expectedGrant, actualGrant)
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

type grantProviderMock struct {
	grant domain.Grant
}

func (r *grantProviderMock) Provide() domain.Grant {
	return r.grant
}
