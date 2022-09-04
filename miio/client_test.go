package miio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRequestDecodeResponse(t *testing.T) {
	// Create a request and then verify that if we decode it
	// we receive the original payload.

	var (
		token   = "0123456789abcdef0123456789abcdef"
		payload = []byte(`{"id": 1, "method": "miIO.info", "params": []}`)
	)

	c, err := New(token)
	require.NoError(t, err)

	req, err := c.createRequest(uint32(1), uint32(2), payload)
	require.NoError(t, err)

	rspPayload, err := c.decodeResponse(req)
	require.NoError(t, err)
	assert.Equal(t, string(payload), string(rspPayload))
}
