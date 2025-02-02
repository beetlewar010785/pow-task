package adapter

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	address             string
	logger              domain.Logger
	listener            net.Listener
	connections         sync.Map
	verifierFactory     application.VerifierFactory
	verificationTimeout time.Duration
}

func StartTCPServer(
	serverAddress string,
	grantProvider domain.QuoteProvider,
	challengeDifficulty domain.Difficulty,
	challengeLength int,
	verificationTimeout time.Duration,
	logger domain.Logger,
) *TCPServer {
	challengeRandomizer := domain.NewSimpleChallengeRandomizer(challengeLength)
	challengeVerifier := domain.NewSimpleChallengeVerifier()
	verifierFactory := application.NewPOWVerifierFactory(
		challengeRandomizer,
		challengeVerifier,
		grantProvider,
		challengeDifficulty,
	)
	return &TCPServer{
		logger:              logger,
		verifierFactory:     verifierFactory,
		address:             serverAddress,
		verificationTimeout: verificationTimeout,
	}
}

func (r *TCPServer) Address() string {
	if r.listener == nil {
		return ""
	}
	return r.listener.Addr().String()
}

func (r *TCPServer) Listen() error {
	listener, err := net.Listen("tcp", r.address)
	if err != nil {
		return fmt.Errorf("failed to start tcp listener on %s: %w", r.address, err)
	}

	r.listener = listener
	r.logger.Debug(fmt.Sprintf("listening on %s", r.address))
	return nil
}

func (r *TCPServer) Run(ctx context.Context) error {
	if r.listener == nil {
		return fmt.Errorf("server is not listening")
	}

	errCh := make(chan error, 1)

	go func() {
		<-ctx.Done()
		_ = r.listener.Close()
		r.closeAllConnections()
		errCh <- ctx.Err()
	}()

	for {
		conn, err := r.listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				r.logger.Warn(fmt.Sprintf("failed to accept connection: %v", err))
				continue
			}
		}

		r.logger.Info(fmt.Sprintf("client connected: %s", conn.RemoteAddr()))

		r.connections.Store(conn, conn)
		go func() {
			defer r.closeConnection(conn)

			err := r.performVerificationWithTimeout(ctx, conn)
			if err != nil {
				r.logger.Warn(fmt.Sprintf("challenge failed: %v", err))
			} else {
				r.logger.Info(fmt.Sprintf("client verified: %s", conn.RemoteAddr()))
			}
		}()
	}
}

func (r *TCPServer) performVerificationWithTimeout(ctx context.Context, conn net.Conn) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, r.verificationTimeout)
	defer cancel()

	return r.performVerification(ctxWithTimeout, conn)
}

func (r *TCPServer) performVerification(ctx context.Context, conn net.Conn) error {
	done := make(chan error, 1)

	go func() {
		readWriter := NewStringReadWriter(conn)
		verifier := r.verifierFactory.Create(readWriter)
		done <- verifier.Verify()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("verification timed out or canceled: %w", ctx.Err())
	case err := <-done:
		return err
	}
}

func (r *TCPServer) closeConnection(conn net.Conn) {
	r.connections.Delete(conn)
	_ = conn.Close()
}

func (r *TCPServer) closeAllConnections() {
	r.connections.Range(func(key, value interface{}) bool {
		conn := value.(net.Conn)
		_ = conn.Close()
		return true
	})
}
