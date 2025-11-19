package header

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const CLRF = "\r\n"
var tokenRe = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+.\-^_` + "`" + `|~]+$`)

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		clrfIndex := strings.Index(string(data), CLRF)
		if clrfIndex == -1 {
			break
		}
		if clrfIndex == 0 {
			done = true
			read += len(CLRF)
			break
		}
		lines := string(data)[:clrfIndex]
		colonIndex := strings.Index(lines, ": ")
		if colonIndex == -1 {
			fmt.Println(lines)
			return 0, false, errors.Join(errors.New("invalid header syntax, no separation of field name and field value"))
		}
		if lines[colonIndex-1] == ' ' {
			return 0, false, errors.Join(errors.New("whitespace not allowed in the field-name"))
		}
		fieldName := strings.ToLower(strings.TrimSpace(lines[:colonIndex]))
		if !isValidToken(fieldName) {
			return 0, false, errors.Join(errors.New("invalid character in the field-name"))
		}
		fieldValue := strings.TrimSpace(lines[colonIndex+2:])
		if v := h.Get(fieldName); v != "" {
			h.Set(fieldName, fmt.Sprintf("%v,%v", v, fieldValue))
		} else {
			h.Set(fieldName, fieldValue)
		}
		consumed := clrfIndex + len(CLRF)
		read += consumed
		data = data[consumed:]
	}
	return read, done, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

func isValidToken(s string) bool {
	return tokenRe.MatchString(s)
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}
