package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/idir-44/httpfromtcp/internal/handler"
	"github.com/idir-44/httpfromtcp/internal/request"
	"github.com/idir-44/httpfromtcp/internal/response"
	"github.com/idir-44/httpfromtcp/internal/server"
)

const port = 42069

func handleRequest(w io.Writer, req *request.Request) *handler.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &handler.HandlerError{
			Code:    response.HTTPStatusBadRequest,
			Message: []byte("Your problem is not my problem\n"),
		}
	case "/myproblem":
		return &handler.HandlerError{
			Code:    response.HTTPStatusInternalServerError,
			Message: []byte("Woopsie, my bad\n"),
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}

}

func main() {
	server, err := server.Serve(port, handleRequest)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
