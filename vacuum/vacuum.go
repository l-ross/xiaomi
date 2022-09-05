package vacuum

import (
	"errors"
	"time"

	"github.com/l-ross/xiaomi/miio"
)

var (
	ErrUnexpectedResponse = errors.New("unexpected response")
)

type Vacuum struct {
	client iClient

	id int64
}

type iClient interface {
	Send(payload []byte) ([]byte, error)
}

func New(c *miio.Client) (*Vacuum, error) {
	v := &Vacuum{
		client: c,
		id:     time.Now().Unix(),
	}

	return v, nil
}

type mockClient struct {
	rsp []byte
	err error
}

func (c *mockClient) Send(payload []byte) ([]byte, error) {
	// If the mock doesn't define a rsp then just return a default.
	if c.rsp == nil {
		c.rsp = []byte(`
			{
				"result": ["ok"],
				"id": 1
			}`,
		)
	}

	return c.rsp, c.err
}
