package adapter

import (
	"context"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/application"
	"github.com/beetlewar010785/pow-task/pkg/lib"
	"net"
	"sync"
)

type TCPServer struct {
	address          string
	logger           lib.Logger
	listener         net.Listener
	connections      sync.Map
	powServerFactory application.POWChallengeHandlerFactory
}

func NewTCPServer(
	address string,
	powServerFactory application.POWChallengeHandlerFactory,
	logger lib.Logger,
) *TCPServer {
	return &TCPServer{
		logger:           logger,
		powServerFactory: powServerFactory,
		address:          address,
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

		in := make(chan []byte)
		out := make(chan []byte)
		tcpConnection := NewTCPConnection(conn, in, out, r.logger)

		r.connections.Store(conn, tcpConnection)

		go r.handleConnection(tcpConnection, in, out)
		go r.runReader(tcpConnection)
		go r.runWriter(tcpConnection)
	}
}

func (r *TCPServer) handleConnection(
	conn *TCPConnection,
	in chan []byte,
	out chan []byte,
) {
	powServer := r.powServerFactory.Create(in, out)
	if err := powServer.Handle(); err != nil {
		r.logger.Warn(fmt.Sprintf("connection handler failed for %s: %v", conn.conn.RemoteAddr(), err))
	}
}

func (r *TCPServer) runReader(conn *TCPConnection) {
	if err := conn.RunReader(); err != nil {
		r.logger.Error(fmt.Sprintf("failed to read from %s: %v", conn.conn.RemoteAddr(), err))
	}
}

func (r *TCPServer) runWriter(conn *TCPConnection) {
	if err := conn.RunWriter(); err != nil {
		r.logger.Error(fmt.Sprintf("failed to write to %s: %v", conn.conn.RemoteAddr(), err))
	}
}

func (r *TCPServer) closeAllConnections() {
	r.connections.Range(func(key, value interface{}) bool {
		conn := value.(*TCPConnection)
		_ = conn.Close()
		return true
	})
}
