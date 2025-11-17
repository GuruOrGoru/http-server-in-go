package header

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFunc(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header with whitespaces
	headers = NewHeaders()
	data = []byte("                Host: localhost:686886               \r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:686886", headers["host"])
	assert.Equal(t, 55, n)
	assert.False(t, done)

	// Test: The final boss
	headers = NewHeaders()
	data = []byte("   Name: guruorgoru   \r\n    Hero: Yes \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 42, n)
	assert.Equal(t, "guruorgoru", headers["name"])
	assert.Equal(t, "Yes", headers["hero"])
	assert.True(t, done)

	// Test: The final boss 2.0
	headers = NewHeaders()
	data = []byte("   N@me: guruorgoru   \r\n    er]o: Yes \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple value boss
	headers = NewHeaders()
	data = []byte(" nAme: guruorgoru \r\n Name: nama \r\n name: wtf \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 49, n)
	assert.Equal(t, "guruorgoru,nama,wtf", headers["name"])
	assert.True(t, done)

}
