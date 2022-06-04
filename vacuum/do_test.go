package vacuum

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVacuum_DoSimple(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"result": ["ok"],
					"id": 12345
				}`,
			),
		},
	}

	err := v.doSimple("test")
	require.NoError(t, err)
}
