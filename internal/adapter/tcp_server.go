package adapter

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"net"
	"sync"
)

type TCPServer struct {
	address           string
	logger            domain.Logger
	listener          net.Listener
	connections       sync.Map
	challengerFactory application.ChallengerFactory
}

func NewTCPServer(
	address string,
	challengerFactory application.ChallengerFactory,
	logger domain.Logger,
) *TCPServer {
	return &TCPServer{
		logger:            logger,
		challengerFactory: challengerFactory,
		address:           address,
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
		go r.performChallenge(conn)
	}
}

func (r *TCPServer) performChallenge(conn net.Conn) {
	defer r.closeConnection(conn)

	challenger := r.challengerFactory.Create(conn, conn)
	if err := challenger.Challenge(); err != nil {
		r.logger.Warn(fmt.Sprintf("challenge for %s failed: %v", conn.RemoteAddr(), err))
	} else {
		r.logger.Debug(fmt.Sprintf("challenge for %s succeeded", conn.RemoteAddr()))
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
