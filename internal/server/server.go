package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/idir-44/httpfromtcp/internal/handler"
	"github.com/idir-44/httpfromtcp/internal/request"
	"github.com/idir-44/httpfromtcp/internal/response"
)

type Server struct {
	handler        handler.Handle
	listener       net.Listener
	isServerClosed atomic.Bool
}

func Serve(port int, handlerFunc handler.Handle) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("error listening to port %d: %v", port, err)
	}
	s := Server{listener: listener, handler: handlerFunc}

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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := handler.HandlerError{
			Code:    response.HTTPStatusBadRequest,
			Message: []byte(err.Error()),
		}
		handler.WriteHandlerError(conn, &handlerErr)
		return
	}
	log.Println(req.RequestLine.Method, req.RequestLine.RequestTarget)
	buff := bytes.NewBuffer([]byte{})

	handlerErr := s.handler(buff, req)
	if handlerErr != nil {
		handler.WriteHandlerError(conn, handlerErr)
		return
	}

	response.WriteStatusLine(conn, response.HTTPStatusOK)
	if err := response.WriteHeaders(conn, response.GetDefaultHeaders(buff.Len())); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	conn.Write(buff.Bytes())

	return
}
