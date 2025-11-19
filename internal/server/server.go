package server

import (
	"bufio"
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

	writer := &response.ResponseWriter{
		Writer: conn,
	}
	s.handler(writer, request)
}

func (s *Server) Close() error {
	s.closed.Store(true)

	if err := s.listener.Close(); err != nil {
		log.Println("Hi, you got an error closing the listener")
	}
	return nil
}
