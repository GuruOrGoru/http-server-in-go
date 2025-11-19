package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/guruorgoru/http-server-in-go/internal/request"
	"github.com/guruorgoru/http-server-in-go/internal/response"
	"github.com/guruorgoru/http-server-in-go/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	checkErr(err)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("Error: Port is not set in the env file")
	}

	server, err := server.Serve(port, handleRoutes)
	checkErr(err)
	defer server.Close()

	log.Println("Server started at port:", port)
	interruptChanel := make(chan os.Signal, 1)
	signal.Notify(interruptChanel, syscall.SIGINT, syscall.SIGTERM)

	<-interruptChanel

	log.Println("Server stopped gracefully")
}

func handleRoutes(w *response.ResponseWriter, req *request.Request) {
	log.Printf("Incoming request: %v %v", req.RequestLine.Method, req.RequestLine.Target)
	switch strings.TrimSpace(req.RequestLine.Target) {
	case "/":
		msg := "GuruOrGoru"
		w.WriteStatusHeader(response.StatusOkay)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	case "/skillissues":
		msg := "You got some skill issues"
		w.WriteStatusHeader(response.StatusBadReq)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	case "/myissues":
		msg := "Sorry, I got some skill issues"
		w.WriteStatusHeader(response.StatusInternalServerError)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	default:
		msg := "No route found in this server :("
		w.WriteStatusHeader(response.StatusNotFound)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error in the server run:", err)
	}
}
