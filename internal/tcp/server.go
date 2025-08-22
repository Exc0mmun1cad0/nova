package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	l "nova/pkg/logger"
	"sync"

	"go.uber.org/zap"
)

var (
	buffSize = 512
)

type Handler interface {
	Serve(context.Context, []byte) []byte
}

type Server struct {
	log *zap.Logger

	// mutex for safe concurrent access from multiple simultanious connections
	mu             sync.Mutex
	connCounter    uint64
	requestCounter uint64

	Addr    string
	Handler Handler
}

func NewServer(addr string, handler Handler, log *zap.Logger) (*Server, error) {
	if log == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Server{
		Addr:           addr,
		Handler:        handler,
		log:            log,
		connCounter:    0,
		requestCounter: 0,
	}, nil
}

func (s *Server) ListenAndServe() {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	s.log.Info("created tcp socket listener", zap.String("address", s.Addr))

	s.log.Info("listening for incoming connections")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err.Error())
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	s.mu.Lock()
	log := s.log.With(zap.Uint64("conn_id", s.connCounter))
	s.connCounter++
	s.mu.Unlock()

	log.Info("accepted new connection")
	for {
		buff := make([]byte, buffSize)
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Error("failed to read request", zap.Error(err))
			continue
		}

		s.mu.Lock()
		log := log.With(zap.Uint64("request_id", s.requestCounter))
		s.requestCounter++
		s.mu.Unlock()

		ctx := l.WithLogger(context.Background(), log)
		resp := s.Handler.Serve(ctx, buff[:n])

		n, err = conn.Write(resp)
		if err != nil {
			log.Error("failed to send response", zap.Error(err))
			continue
		}
		log.Info("sent response", zap.Int("bytes", n))
	}

	log.Info("connection closed")
}
