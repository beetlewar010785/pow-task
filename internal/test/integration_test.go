package test

import (
	"context"
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
		server            *adapter.POWServer
		client            net.Conn
		solver            application.Solver
		ctx               context.Context
	}

	setup := func(
		t *testing.T,
		challengeDifficulty domain.Difficulty,
		challengeLength int,
	) testSuite {
		const randomPort = ":0"

		grantProvider := new(grantProviderMock)
		powServer := adapter.StartPOWServer(
			randomPort,
			grantProvider,
			challengeDifficulty,
			challengeLength,
			time.Second,
			adapter.NewStdLogger("server", adapter.LogLevelInfo),
		)

		require.NoError(t, powServer.Listen())

		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		go func() {
			err := powServer.Run(ctx)
			assert.NoError(t, err, "server.Run() returned an unexpected error")
		}()

		WaitForServer(t, powServer)

		conn, solver, err := adapter.CreatePOWClient(powServer.Address(), 10*time.Second)
		require.NoError(t, err)

		return testSuite{
			grantProvider,
			powServer,
			conn,
			solver,
			ctx,
		}
	}

	tearDown := func(suite testSuite) {
		suite.ctx.Done()
		_ = suite.client.Close()
	}

	t.Run("should pass POW challenge and get grant", func(t *testing.T) {
		suite := setup(t, 2, 10)
		defer tearDown(suite)

		suite.grantProviderMock.grant = "expected-grant"
		actualGrant, err := suite.solver.Solve()
		require.NoError(t, err)

		expectedGrant := domain.SuccessGrant(suite.grantProviderMock.grant)
		assert.Equal(t, expectedGrant, actualGrant)
	})
}

func WaitForServer(t *testing.T, server *adapter.POWServer) {
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
	grant domain.Quote
}

func (r *grantProviderMock) Provide() domain.Quote {
	return r.grant
}
