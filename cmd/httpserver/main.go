package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/idir-44/httpfromtcp/internal/request"
	"github.com/idir-44/httpfromtcp/internal/response"
	"github.com/idir-44/httpfromtcp/internal/server"
)

const port = 42069

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

func handleRequest(res *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		res.WriteStatusLine(response.HTTPStatusBadRequest)
		body := []byte(`
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Override("Content-type", "text/html")
		res.WriteHeaders(headers)

		res.WriteBody(body)
	case "/myproblem":
		res.WriteStatusLine(response.HTTPStatusInternalServerError)
		body := []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Override("Content-type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(body)
	default:
		res.WriteStatusLine(response.HTTPStatusOK)
		body := []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
		`)
		headers := response.GetDefaultHeaders(len(body))
		headers.Override("Content-type", "text/html")
		res.WriteHeaders(headers)
		res.WriteBody(body)
	}
}
