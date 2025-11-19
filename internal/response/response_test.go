package response

import (
	"bytes"
	"testing"

	"github.com/guruorgoru/http-server-in-go/internal/header"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteStatusHeader(t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteStatusHeader(buf, StatusOkay)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 200 OK\r\n", buf.String())

	buf.Reset()
	err = WriteStatusHeader(buf, StatusBadReq)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 400 Bad Request\r\n", buf.String())

	buf.Reset()
	err = WriteStatusHeader(buf, StatusInternalServerError)
	require.NoError(t, err)
	assert.Equal(t, "HTTP/1.1 500 Internal Server Error\r\n", buf.String())
}

func TestGetDefaultHeaders(t *testing.T) {
	headers := GetDefaultHeaders(10)
	assert.Equal(t, "10", headers.Get("content-length"))
	assert.Equal(t, "close", headers.Get("connection"))
	assert.Equal(t, "text/plain", headers.Get("content-type"))
}

func TestWriteHeaders(t *testing.T) {
	buf := &bytes.Buffer{}
	headers := header.NewHeaders()
	headers.Set("test", "value")
	headers.Set("another", "header")
	err := WriteHeaders(buf, headers)
	require.NoError(t, err)
	// Order may vary, check contains
	output := buf.String()
	assert.Contains(t, output, "test: value\r\n")
	assert.Contains(t, output, "another: header\r\n")
	assert.Contains(t, output, "\r\n")
}

func TestWriteError(t *testing.T) {
	buf := &bytes.Buffer{}
	hErr := &HandlingError{StatusCode: StatusBadReq, Msg: "bad request"}
	err := WriteError(buf, hErr)
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "HTTP/1.1 400 Bad Request\r\n")
	assert.Contains(t, output, "content-length:")
	assert.Contains(t, output, "connection: close")
	assert.Contains(t, output, "content-type: text/plain")
	assert.Contains(t, output, "\r\nError:\n- Code: 400\n- Message: bad request")
}
