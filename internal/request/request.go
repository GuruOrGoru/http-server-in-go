package request

import (
	"errors"
	"io"
	"strings"
)

type parsingState int

const (
	InitState parsingState = iota
	DoneState
)

type Request struct {
	RequestLine RequestLine
	Header      map[string]string
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
		State: InitState,
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
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == DoneState {
		return 0, nil
	}
	requestLine, bytesRead, err := requestLineParse(data)
	if err != nil {
		return 0, err
	}
	if requestLine != nil {
		r.RequestLine = *requestLine
		r.State = DoneState
	}
	return bytesRead, err
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
