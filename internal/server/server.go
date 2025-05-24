package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/idir-44/httpfromtcp/internal/request"
	"github.com/idir-44/httpfromtcp/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message []byte
}

type Handle func(res *response.Writer, req *request.Request)

func WriteHandlerError(w io.Writer, errHandler *HandlerError) error {
	err := response.WriteSatusLine(w, errHandler.Code)
	if err != nil {
		return err
	}

	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(errHandler.Message)))
	if err != nil {
		return err
	}

	_, err = w.Write(errHandler.Message)

	return err
}

type Server struct {
	handler        Handle
	listener       net.Listener
	isServerClosed atomic.Bool
}

func Serve(port int, handlerFunc Handle) (*Server, error) {
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
	res := response.NewReponseWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println(req.RequestLine.Method, req.RequestLine.RequestTarget, response.HTTPStatusBadRequest)
		res.WriteStatusLine(response.HTTPStatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing the request: %v", err))
		res.WriteHeaders(response.GetDefaultHeaders(len(body)))
		res.WriteBody(body)
		return
	}

	s.handler(res, req)
	log.Println(req.RequestLine.Method, req.RequestLine.RequestTarget)

	return
}
