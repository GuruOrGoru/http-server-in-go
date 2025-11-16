package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
		lines := getLinesChannel(reader, conn)
		for line := range lines {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(reader *bufio.Reader, f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatalln("error closing file: ", err)
			}
		}()
		defer close(ch)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					ch <- line
					break
				}
				log.Fatalln("error reading from messages.txt file", err)
			}
			ch <-line
		}
	}()
	return ch
}
