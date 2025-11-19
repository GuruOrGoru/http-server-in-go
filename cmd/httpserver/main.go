package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
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
	switch {
	case strings.TrimSpace(req.RequestLine.Target) == "/":
		msg := "GuruOrGoru"
		w.WriteStatusHeader(response.StatusOkay)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	case strings.TrimSpace(req.RequestLine.Target) == "/skillissues":
		msg := "You got some skill issues"
		w.WriteStatusHeader(response.StatusBadReq)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	case strings.TrimSpace(req.RequestLine.Target) == "/myissues":
		msg := "Sorry, I got some skill issues"
		w.WriteStatusHeader(response.StatusInternalServerError)
		w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
		w.Writer.Write([]byte(msg))
	case strings.HasPrefix(req.RequestLine.Target, "/httpbin/stream/"):
		respy, err := http.Get("https://httpbin.org/" + req.RequestLine.Target[len("/httpbin/"):])
		respyBody := []byte{}
		if err != nil {
			msg := "Sorry, I got some skill issues"
			w.WriteStatusHeader(response.StatusInternalServerError)
			w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
			w.Writer.Write([]byte(msg))
		}
		w.WriteStatusHeader(response.StatusOkay)
		headers := response.GetDefaultHeaders(0)
		headers.Delete("content-length")
		headers.Set("transfer-encoding", "chunked")
		headers.Set("Trailers", "X-Content-SHA256")
		headers.Set("Trailers", "X-Content-Length")
		w.WriteHeaders(headers)
		for {
			data := make([]byte, 1024)
			n, err := respy.Body.Read(data)
			if err != nil {
				break
			}
			respyBody = append(respyBody, data[:n]...)
			w.Writer.Write([]byte(fmt.Sprintf("%x\r\n", n)))
			w.Writer.Write(data[:n])
			w.Writer.Write([]byte("\r\n"))
		}
		w.Writer.Write([]byte("0\r\n"))
		trailers := response.GetDefaultHeaders(len(respyBody))
		trailers.Delete("content-length")
		hash := sha256.Sum256(respyBody)
		trailers.Set("X-Content-SHA256", hex.EncodeToString(hash[:]))
		trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(respyBody)))
		w.WriteTrailers(trailers)
	case strings.TrimSpace(req.RequestLine.Target) == "/video":
		f, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			msg := "Sorry, I got some skill issues"
			w.WriteStatusHeader(response.StatusInternalServerError)
			w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
			w.Writer.Write([]byte(msg))
		}
		w.WriteStatusHeader(response.StatusOkay)
		h := response.GetDefaultHeaders(len(f))
		h.Set("content-type", "video/mp4")
		w.WriteHeaders(h)
		w.Writer.Write(f)

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
