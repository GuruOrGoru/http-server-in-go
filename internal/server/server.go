package server

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"sync/atomic"

	"github.com/guruorgoru/http-server-in-go/internal/request"
	"github.com/guruorgoru/http-server-in-go/internal/response"
)

type Server struct {
	listener net.Listener
	handler  response.Handler
	closed   atomic.Bool
}

func Serve(port string, handler response.Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
		handler:  handler,
	}
	s.closed.Store(false)
	go func() {
		s.listen()
	}()
	return s, err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Println("Error accepting the connection")
			continue
		}
		go func() {
			s.handle(conn)
		}()
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	request, err := request.RequestFromReader(reader)
	if err != nil {
		log.Println("Got an error while reading from the request:", err)
		return
	}

	writeBuf := bytes.NewBuffer([]byte{})
	respBuff := bytes.NewBuffer([]byte{})

	hErr := s.handler(writeBuf, request)
	if hErr != nil {
		if err = response.WriteError(conn, hErr); err != nil {
			log.Println("Error writing handling err:", err)
		}
		return
	}

	if err = response.WriteStatusHeader(respBuff, response.StatusOkay); err != nil {
		log.Println("Error writing statusLine:", err)
		return

	}

	headers := response.GetDefaultHeaders(writeBuf.Len())
	if err = response.WriteHeaders(respBuff, headers); err != nil {
		log.Println("Error writing headers:", err)
		return

	}
	respBuff.Write(writeBuf.Bytes())
	if _, err = conn.Write(respBuff.Bytes()); err != nil {
		log.Println("Error writing body:", err)
		return
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)

	if err := s.listener.Close(); err != nil {
		log.Println("Hi, you got an error closing the listener")
	}
	return nil
}
