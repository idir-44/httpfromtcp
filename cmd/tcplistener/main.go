package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/idir-44/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error listening to address: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error getting request %v", err)
			return
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)
		if len(req.Headers) > 0 {
			fmt.Println("Headers:")
			for key, value := range req.Headers {
				fmt.Printf("- %s: %s\n", key, value)
			}
		}

		if len(req.Body) > 0 {
			fmt.Println("Body:")
			fmt.Println(string(req.Body))

		}

	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func(f io.ReadCloser) {
		defer f.Close()
		defer close(lines)
		currentLine := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				panic(err)
			}

			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s\n", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}(f)

	return lines
}
