//go:build integration

package vacuum

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"

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

	i, err := v.GetDNDTimer()
	require.NoError(t, err)

	spew.Dump(i)
}
