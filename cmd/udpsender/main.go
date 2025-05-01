package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	resolveAddr := "localhost:42069"

	udpAddr, err := net.ResolveUDPAddr("udp", resolveAddr)
	if err != nil {
		log.Fatalf("Couldn't resolve %s", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Couldn't dial: %s", err.Error())
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}
	}

}
