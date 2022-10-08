//go:build integration

package miio

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	// Basic integration test to send an info request to a Xiaomi Robot Vacuum and validate the response.

	ip := os.Getenv("MIIO_IP")
	token := os.Getenv("MIIO_TOKEN")
	// Use the timestamp for the ID to ensure that the test can be run multiple times.
	payload := []byte(fmt.Sprintf(`{"id": %d, "method": "miIO.info", "params": []}`, time.Now().Unix()))

	c, err := New(token, SetIP(ip))
	require.NoError(t, err)

	err = c.Connect()
	require.NoError(t, err)

	rsp, err := c.Send(payload)
	require.NoError(t, err)

	err = c.Close()
	require.NoError(t, err)

	rspBody := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(rsp, &rspBody)
	require.NoError(t, err)
	assert.Equal(t, "ok", rspBody.Message)
}
