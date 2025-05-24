package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/idir-44/httpfromtcp/internal/headers"
)

const crlf = "\r\n"

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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != WriterStateHeader {
		return fmt.Errorf("cannot write headers in state %s", w.state)
	}
	defer func() { w.state = WriterStateBody }()
	return WriteHeaders(w.writer, headers)
}
