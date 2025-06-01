package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		res.WriteStatusLine(response.HTTPStatusOK)
		headers := response.GetDefaultHeaders(0)
		_, ok := headers.Delete("Content-Length")
		if !ok {
			log.Println("key Content-Length not found")
		}
		headers.Override("Connection", "Keep-Alive")
		headers.Set("Transfer-Encoding", "chunked")
		res.WriteHeaders(headers)
		err := proxyChunkedDate(res, strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin"))
		if err != nil {
			log.Println("error reading chunks: ", err)
		}
		return
	}

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

func proxyChunkedDate(res *response.Writer, route string) error {
	response, err := http.Get(fmt.Sprintf("https://httpbin.org%s", route))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buff := make([]byte, 1024)

	for {
		n, err := response.Body.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				res.WriteChunkedBodyDone()
				return nil
			}
			return err
		}

		_, err = res.WriteChunkedBody([]byte(fmt.Sprintf("%X", n)))
		if err != nil {
			return err
		}
		n, err = res.WriteChunkedBody(buff[:n])
		if err != nil {
			return err
		}
	}
}
