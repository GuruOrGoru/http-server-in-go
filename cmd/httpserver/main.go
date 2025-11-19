package main

import (
	"io"
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

func handleRoutes(w io.Writer, req *request.Request) *response.HandlingError {
	log.Printf("Incoming request: %v %v\n", req.RequestLine.Method, req.RequestLine.Target)
	switch strings.TrimSpace(req.RequestLine.Target) {
	case "/":
		_, err := w.Write([]byte("Hello There!"))
		if err != nil {
			return &response.HandlingError{
				StatusCode: 500,
				Msg:        "error writing data: " + err.Error(),
			}
		}
		return nil

	case "/skillissues":
		return &response.HandlingError{
			StatusCode: 400,
			Msg:        "you got skill issues",
		}
	case "/myissues":
		return &response.HandlingError{
			StatusCode: 500,
			Msg:        "i have skill issues :(",
		}
	default:
		return &response.HandlingError{
			StatusCode: 404,
			Msg:        "page not found",
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln("Error in the server run:", err)
	}
}
