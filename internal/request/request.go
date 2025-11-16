package request

type Request struct {
	RequestLine RequestLine
	Header      map[string]string
	Body        []byte
}

type RequestLine struct {
	Method      string
	Target      string
	HttpVersion string
}

func RequestLineParse() {
}
