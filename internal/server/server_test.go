package server

import (
	"io"
	"net"
	"testing"

	"github.com/guruorgoru/http-server-in-go/internal/request"
	"github.com/guruorgoru/http-server-in-go/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeInvalidPort(t *testing.T) {
	_, err := Serve("invalid", nil)
	require.Error(t, err)
}

func TestServeValidPort(t *testing.T) {
	s, err := Serve("0", func(w io.Writer, req *request.Request) *response.HandlingError {
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	assert.NotNil(t, s.listener)
	assert.False(t, s.closed.Load())
	s.Close()
}

func TestClose(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	defer listener.Close()

	s := &Server{listener: listener}
	err = s.Close()
	assert.NoError(t, err)
	assert.True(t, s.closed.Load())
}
