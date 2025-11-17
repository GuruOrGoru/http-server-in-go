package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/guruorgoru/http-server-in-go/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", "42069"))
	if err != nil {
		log.Fatalln("error making a new tcp listener", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("error accepting the conenction", err)
		}
		reader := bufio.NewReader(conn)
		requestData, err := request.RequestFromReader(reader)
		if err != nil {
			log.Fatalln("error while getting the request data from the reader:", err)
		}
		textToPrint := fmt.Sprintf("Request line:\n- Method: %v\n- Target: %v\n- Version: %v", requestData.RequestLine.Method, requestData.RequestLine.Target, requestData.RequestLine.HttpVersion)
		fmt.Println(textToPrint)
	}
}
