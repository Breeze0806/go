package pqto

import (
	"context"
	"net"
	"time"

	"github.com/Breeze0806/go/timeout"
)

type dialer struct {
	*timeout.Dialer
}

func newDialer(readTimeout time.Duration, writeTimeout time.Duration) *dialer {
	return &dialer{
		Dialer: timeout.NewDialer(readTimeout, writeTimeout),
	}
}

func (d *dialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return d.DialContext(ctx, network, address)
}
