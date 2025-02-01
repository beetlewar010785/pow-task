package adapter

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"net"
)

type TCPClient struct {
	serverAddress string
	in            chan []byte
	out           chan []byte
	logger        domain.Logger
}

func NewTCPClient(
	serverAddress string,
	in chan []byte,
	out chan []byte,
	logger domain.Logger,
) *TCPClient {
	return &TCPClient{
		serverAddress,
		in,
		out,
		logger,
	}
}

func (r *TCPClient) Run(ctx context.Context) error {
	conn, err := net.Dial("tcp", r.serverAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", r.serverAddress, err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	connection := NewTCPConnection(
		conn,
		r.in,
		r.out,
		r.logger,
	)

	errCh := make(chan error, 2)

	go func() {
		if err := connection.RunReader(); err != nil {
			r.logger.Error(fmt.Sprintf("reader error: %v", err))
			errCh <- err
		}
	}()

	go func() {
		if err := connection.RunWriter(); err != nil {
			r.logger.Error(fmt.Sprintf("writer error: %v", err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return fmt.Errorf("connection error: %w", err)
	}
}
