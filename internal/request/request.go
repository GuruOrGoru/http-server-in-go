package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/guruorgoru/http-server-in-go/internal/header"
)

type parsingState int

const (
	InitState parsingState = iota
	DoneState
	HeaderState
	BodyState
)

type Request struct {
	RequestLine RequestLine
	Headers     header.Headers
	Body        []byte
	State       parsingState
}

type RequestLine struct {
	Method      string
	Target      string
	HttpVersion string
}

const (
	REQUIRED_HTTP_VERSION string = "HTTP/1.1"
	CLRF                  string = "\r\n"
)

func RequestFromReader(r io.Reader) (*Request, error) {
	request := &Request{
		State:   InitState,
		Headers: header.NewHeaders(),
	}
	temporaryBuffer := make([]byte, 8)
	internalBuffer := make([]byte, 0, 8)
	for request.State != DoneState {
		n, err := r.Read(temporaryBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Join(errors.New("Error reading to a buffer:"), err)
		}
		if n > 0 {
			internalBuffer = append(internalBuffer, temporaryBuffer[:n]...)
		}
		consumedByte, err := request.parse(internalBuffer)
		if err != nil {
			return nil, err
		}
		internalBuffer = internalBuffer[consumedByte:]
	}
	if val, ok := request.Headers["content-length"]; ok {
		expcLength, _ := strconv.Atoi(val)
		if len(request.Body) != expcLength {
			return nil, errors.Join(fmt.Errorf("brutha, your content-length isnt equal to the body length(i.e %v != %v)", expcLength, len(request.Body)))
		}
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case InitState:
		requestLine, bytesRead, err := requestLineParse(data)
		if err != nil {
			return 0, err
		}
		if requestLine != nil {
			r.RequestLine = *requestLine
			r.State = HeaderState
		}
		return bytesRead, err
	case DoneState:
		return 0, nil
	case HeaderState:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = BodyState
		}
		return n, nil
	case BodyState:
		v, ok := r.Headers["content-length"]
		if !ok || v == "" {
			r.State = DoneState
			return 0, nil
		}
		iv, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		remaining := iv - len(r.Body)
		if remaining <= 0 {
			r.State = DoneState
			return 0, nil
		}
		toRead := min(len(data), remaining)
		fullBody := data[:toRead]
		r.Body = append(r.Body, fullBody...)
		if len(r.Body) == iv {
			r.State = DoneState
		}
		return toRead, nil
	default:
		panic("bruh your code sucks")
	}
}

func requestLineParse(content []byte) (*RequestLine, int, error) {
	clrfIndex := strings.Index(string(content), CLRF)
	if clrfIndex == -1 {
		return nil, 0, nil
	}
	lines := string(content)[:clrfIndex]
	nextContentIndex := clrfIndex + 2
	requestLineArgs := strings.Split(lines, " ")
	if len(requestLineArgs) != 3 {
		return nil, nextContentIndex, errors.Join(errors.New("not enough request arguments in request line:"), errors.New(lines))
	}
	methodArg := requestLineArgs[0]
	if !checkAlphabetsAndCapitalization(methodArg) {
		return nil, nextContentIndex, errors.New("method is either not provided or is incorrect")
	}
	targetArg := requestLineArgs[1]
	httpVersion := requestLineArgs[2]
	if !checkHttpVersion(httpVersion) {
		return nil, nextContentIndex, errors.New("only http version 1.1 is supported")
	}
	return &RequestLine{
		Method:      methodArg,
		Target:      targetArg,
		HttpVersion: httpVersion[5:],
	}, nextContentIndex, nil
}

func checkHttpVersion(httpVersion string) bool {
	return httpVersion == REQUIRED_HTTP_VERSION
}

func checkAlphabetsAndCapitalization(str string) bool {
	totalCaps := 0
	for _, ch := range str {
		if ch >= 'A' && ch <= 'Z' {
			totalCaps++
		}
	}
	return totalCaps == len(str)
}
