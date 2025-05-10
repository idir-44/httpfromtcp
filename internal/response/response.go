package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/idir-44/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	HTTPStatusOK                  StatusCode = 200
	HTTPStatusBadRequest          StatusCode = 400
	HTTPStatusInternalServerError StatusCode = 500
)

const crlf = "\r\n"

// TODO: move status line code to it's own file
func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error

	//TODO: refactor that
	switch statusCode {
	case HTTPStatusOK:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d OK%s", statusCode, crlf)))
	case HTTPStatusBadRequest:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Bad Request%s", statusCode, crlf)))
	case HTTPStatusInternalServerError:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Internal Server Error%s", statusCode, crlf)))
	default:
		_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s", statusCode, crlf)))
	}

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	h.Set("Content-length", strconv.Itoa(contentLen))

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s%s", key, value, crlf)))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
