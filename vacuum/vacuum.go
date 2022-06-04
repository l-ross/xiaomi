package vacuum

import (
	"errors"
	"time"

	"github.com/l-ross/miio/client"
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

func New(c *client.Client) (*Vacuum, error) {
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
	return c.rsp, c.err
}
