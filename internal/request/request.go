package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/idir-44/httpfromtcp/internal/headers"
)

type ParserState int

const (
	ParserStateInitialized ParserState = iota
	ParserStateParsingHeaders
	ParserStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
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
	request.Headers = headers.NewHeaders()

	for request.parserState != ParserStateDone {
		nbBytesRead, err := reader.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.parserState != ParserStateDone {
					return nil, fmt.Errorf("incomplete request, in state %d, read n bytes on EOF: %d", request.parserState, nbBytesRead)
				}
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

	totalBytesParsed := 0
	for r.parserState != ParserStateDone {
		nbParsedBytes, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += nbParsedBytes
		if nbParsedBytes == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}
func (r *Request) parseSingle(data []byte) (int, error) {
	if r.parserState == ParserStateDone {
		return 0, fmt.Errorf("error: cannot read data in a done state")
	}

	switch r.parserState {
	case ParserStateInitialized:
		nbParsedLineBytes, err := r.parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		return nbParsedLineBytes, nil
	case ParserStateParsingHeaders:
		nbParsedHeaderBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.parserState = ParserStateDone
		}
		return nbParsedHeaderBytes, nil
	default:
		return 0, fmt.Errorf("unknown parser state")
	}
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	requestLine, nbParsedBytes, err := parseRequestLine(data)
	if err != nil {
		return nbParsedBytes, err
	}

	if nbParsedBytes == 0 {
		return nbParsedBytes, nil
	}

	r.parserState = ParserStateParsingHeaders
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
