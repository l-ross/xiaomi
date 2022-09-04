//go:build integration

package vacuum

import (
	"os"
	"testing"

	"github.com/l-ross/xiaomi/miio"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ip := os.Getenv("MIIO_IP")
	token := os.Getenv("MIIO_TOKEN")

	c, err := miio.New(token, miio.SetIP(ip))
	require.NoError(t, err)

	v, err := New(c)
	require.NoError(t, err)

	err = v.StartCharge()
	require.NoError(t, err)
}
