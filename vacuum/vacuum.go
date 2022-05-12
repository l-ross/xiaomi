package vacuum

import (
	"time"

	"github.com/l-ross/miio/client"
)

type Vacuum struct {
	client *client.Client

	id int64
}

func New(c *client.Client) (*Vacuum, error) {
	v := &Vacuum{
		client: c,
		id:     time.Now().Unix(),
	}

	return v, nil
}
