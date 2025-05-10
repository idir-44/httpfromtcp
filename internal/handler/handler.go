package handler

import (
	"io"

	"github.com/idir-44/httpfromtcp/internal/request"
	"github.com/idir-44/httpfromtcp/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message []byte
}

type Handle func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, errHandler *HandlerError) error {
	err := response.WriteStatusLine(w, errHandler.Code)
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
