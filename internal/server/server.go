package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/idir-44/httpfromtcp/internal/response"
)

type Server struct {
	listener       net.Listener
	isServerClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error listening to port %d: %v", port, err)
	}
	s := Server{listener: listener}

	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	if s.isServerClosed.Load() {
		return fmt.Errorf("server already closed")
	}

	s.isServerClosed.Store(true)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isServerClosed.Load() {
				break
			}
			log.Fatalf("error accepting connection: %v", err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	response.WriteStatusLine(conn, response.HTTPStatusOK)
	if err := response.WriteHeaders(conn, response.GetDefaultHeaders(0)); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return
}
