package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type ParserState string

const (
	ParserStateDone        ParserState = "done"
	ParserStateInitialized ParserState = "initialized"
)

type Request struct {
	RequestLine RequestLine
	parserState ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const BUFF_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	var buff []byte
	tmp := make([]byte, BUFF_SIZE, BUFF_SIZE)
	request := &Request{parserState: ParserStateInitialized}

	for request.parserState != ParserStateDone {
		nbBytesRead, err := reader.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.parserState = ParserStateDone
				break
			}
			return nil, err
		}

		if nbBytesRead > 0 {
			buff = append(buff, tmp[:nbBytesRead]...)
		}

		nbParsedData, err := request.parse(buff)
		if err != nil {
			return nil, err
		}

		if nbParsedData > 0 {
			buff = buff[nbParsedData:]
		}

	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.parserState == ParserStateDone {
		return 0, fmt.Errorf("error: cannot read data in a done state")
	}

	if r.parserState != ParserStateDone && r.parserState != ParserStateInitialized {
		return 0, fmt.Errorf("error: unknown parser state")
	}

	requestLine, nbParsedBytes, err := parseRequestLine(data)
	if err != nil {
		return nbParsedBytes, err
	}

	if requestLine == nil {
		return nbParsedBytes, nil
	}

	r.parserState = ParserStateDone
	r.RequestLine = *requestLine

	return nbParsedBytes, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:idx])

	if requestLineText == "" {
		return nil, 0, fmt.Errorf("request should not be empty")
	}

	requestLine := strings.Split(requestLineText, " ")
	if len(requestLine) != 3 {
		return nil, 0, fmt.Errorf("request-line should have 3 parts")
	}

	method := requestLine[0]
	requestTarget := requestLine[1]
	httpVersion := requestLine[2]

	for _, s := range method {
		if !unicode.IsUpper(s) || !unicode.IsLetter(s) {
			return nil, 0, fmt.Errorf("method should be alphabetic charachrers and uppercase")
		}
	}

	if httpVersion != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("http version should be HTTP/1.1")
	}

	return &RequestLine{
		Method:        method,
		HttpVersion:   httpVersion[5:],
		RequestTarget: requestTarget,
	}, idx + 2, nil
}
