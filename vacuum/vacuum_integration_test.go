//go:build integration

package vacuum

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/l-ross/miio/client"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ip := os.Getenv("MIIO_IP")
	token := os.Getenv("MIIO_TOKEN")

	c, err := client.New(token, client.SetIP(ip))
	require.NoError(t, err)

	v, err := New(c)
	require.NoError(t, err)

	//err = v.StartCharge()
	l, err := v.FirmwareFeatures()
	require.NoError(t, err)

	spew.Dump(l)
}
