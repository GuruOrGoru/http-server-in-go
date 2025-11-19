package main

import (
	"bytes"
	"testing"

	"github.com/guruorgoru/http-server-in-go/internal/request"
	"github.com/guruorgoru/http-server-in-go/internal/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleRoutes(t *testing.T) {
	tests := []struct {
		name         string
		target       string
		expectedBody string
		expectedErr  *response.HandlingError
	}{
		{
			name:         "root path",
			target:       "/",
			expectedBody: "Hello There!",
			expectedErr:  nil,
		},
		{
			name:   "skill issues",
			target: "/skillissues",
			expectedErr: &response.HandlingError{
				StatusCode: 400,
				Msg:        "you got skill issues",
			},
		},
		{
			name:   "my issues",
			target: "/myissues",
			expectedErr: &response.HandlingError{
				StatusCode: 500,
				Msg:        "i have skill issues :(",
			},
		},
		{
			name:   "not found",
			target: "/unknown",
			expectedErr: &response.HandlingError{
				StatusCode: 404,
				Msg:        "page not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			req := &request.Request{
				RequestLine: request.RequestLine{
					Target: tt.target,
				},
			}

			err := handleRoutes(&buf, req)

			if tt.expectedErr == nil {
				require.Nil(t, err)
				assert.Equal(t, tt.expectedBody, buf.String())
			} else {
				require.NotNil(t, err)
				assert.Equal(t, tt.expectedErr.StatusCode, err.StatusCode)
				assert.Equal(t, tt.expectedErr.Msg, err.Msg)
			}
		})
	}
}
