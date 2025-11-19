package response

import (
	"fmt"
	"io"

	"github.com/guruorgoru/http-server-in-go/internal/header"
	"github.com/guruorgoru/http-server-in-go/internal/request"
)

type (
	Handler       func(w io.Writer, req *request.Request) *HandlingError
	HandlingError struct {
		StatusCode StatusCode
		Msg        string
	}
)

type StatusCode int

const (
	StatusOkay                StatusCode = 200
	StatusBadReq              StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusNotFound            StatusCode = 404
)

func WriteStatusHeader(w io.Writer, statusCode StatusCode) error {
	reason := ""
	switch statusCode {
	case StatusOkay:
		reason = "OK"
	case StatusBadReq:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "Internal Server Error"
	case StatusNotFound:
		reason = "Page Not Found"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %v %v\r\n", statusCode, reason)
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(conLength int) header.Headers {
	newHeader := header.NewHeaders()
	newHeader.Set("Content-Length", fmt.Sprint(conLength))
	newHeader.Set("Connection", "close")
	newHeader.Set("Content-Type", "text/plain")
	return newHeader
}

func WriteHeaders(w io.Writer, headers header.Headers) error {
	headerMsg := ""
	for name, value := range headers {
		headerMsg += fmt.Sprintf("%v: %v\r\n", name, value)
	}
	headerMsg += "\r\n"
	_, err := w.Write([]byte(headerMsg))
	return err
}

func WriteError(w io.Writer, hErr *HandlingError) error {
	if err := WriteStatusHeader(w, hErr.StatusCode); err != nil {
		return err
	}

	errorMsg := fmt.Sprintf("Error:\n- Code: %v\n- Message: %v", hErr.StatusCode, hErr.Msg)

	headers := GetDefaultHeaders(len(errorMsg))

	if err := WriteHeaders(w, headers); err != nil {
		return err
	}

	_, err := w.Write([]byte(errorMsg))
	return err
}
