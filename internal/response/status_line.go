package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	HTTPStatusOK                  StatusCode = 200
	HTTPStatusBadRequest          StatusCode = 400
	HTTPStatusInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case HTTPStatusOK:
		reasonPhrase = "OK"
	case HTTPStatusBadRequest:
		reasonPhrase = "Bas Request"
	case HTTPStatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteSatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))

	return err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateStatusLine {
		return fmt.Errorf("cannot write status line in state: %s", w.state)
	}

	defer func() { w.state = WriterStateHeader }()

	return WriteSatusLine(w.writer, statusCode)
}
