package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

// TODO: move to configuration module (or package)
var (
	addr = "localhost:6379"
)

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to init socket listener: %v", err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept client connection: %v", err.Error())
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	for {
		buff := make([]byte, 100)
		n, err := conn.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed. Maybe...")
				break
			}

			fmt.Printf("failed to read request due to: %v", err.Error())
		}

		req := buff[:n]
		fmt.Printf("Got request: %s\n", string(req))
	}
}
