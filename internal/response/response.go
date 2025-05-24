package response

import (
	"io"
)

type WriterState = string

const (
	WriterStateStatusLine WriterState = "status_line"
	WriterStateHeader     WriterState = "headers"
	WriterStateBody       WriterState = "body"
)

type Writer struct {
	writer io.Writer
	state  WriterState
}

func NewReponseWriter(w io.Writer) *Writer {
	return &Writer{writer: w, state: WriterStateStatusLine}
}
