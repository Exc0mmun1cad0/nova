package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
)

var (
	buffSize = 512
)

type Handler interface {
	Serve([]byte) []byte
}

type Server struct {
	Addr    string
	Handler Handler
}

func (s *Server) ListenAndServe() {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err.Error())
			continue
		}

		fmt.Println("got new connection")
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	for {
		buff := make([]byte, buffSize)
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed. Maybe...")
				break
			}

			fmt.Printf("failed to read request due to: %v", err.Error())
		}

		// cmd, _ := resp.Decode(buff[:n])
		resp := s.Handler.Serve(buff[:n])
		fmt.Println("response:", string(resp))

		n, err = conn.Write(resp)
		fmt.Println(n, err)
	}
}
