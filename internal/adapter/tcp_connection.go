package adapter

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/beetlewar010785/pow-task/pkg/lib"
	"net"
)

type TCPConnection struct {
	conn   net.Conn
	in     chan []byte
	out    chan []byte
	logger lib.Logger
}

func NewTCPConnection(
	conn net.Conn,
	in chan []byte,
	out chan []byte,
	logger lib.Logger,
) *TCPConnection {
	return &TCPConnection{
		conn,
		in,
		out,
		logger,
	}
}

func (r *TCPConnection) RunReader() error {
	scanner := bufio.NewScanner(r.conn)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		r.logger.Trace(fmt.Sprintf("read bytes: %d", len(bytes)))
		r.in <- bytes
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("connection read error: %w", err)
	}

	return nil
}

func (r *TCPConnection) RunWriter() error {
	for {
		select {
		case bytes := <-r.out:
			bytes = append(bytes, "\n"...)
			num, err := r.conn.Write(bytes)

			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			if err != nil {
				return fmt.Errorf("failed to write bytes: %w", err)
			}

			r.logger.Trace(fmt.Sprintf("written bytes: %d", num))
		}
	}
}

func (r *TCPConnection) Close() error {
	return r.conn.Close()
}
