package timeout

import (
	"context"
	"net"
	"time"
)

type dialer interface {
	dial(network, address string) (net.Conn, error)
	dialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type defaultDialer struct {
	d net.Dialer
}

func (d *defaultDialer) dial(network, address string) (net.Conn, error) {
	return d.d.Dial(network, address)
}

func (d *defaultDialer) dialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.d.DialContext(ctx, network, address)
}

type Dialer struct {
	dialer
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewDialer(readTimeout time.Duration, writeTimeout time.Duration) *Dialer {
	return &Dialer{
		dialer:       &defaultDialer{},
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

func (d *Dialer) Dial(network, address string) (conn net.Conn, err error) {
	conn, err = d.dial(network, address)
	if err != nil {
		return nil, err
	}
	return d.conn(conn)
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	conn, err = d.dialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return d.conn(conn)
}

func (d *Dialer) conn(c net.Conn) (net.Conn, error) {
	if d.readTimeout == 0 && d.writeTimeout == 0 {
		return c, nil
	}
	return NewConn(c, d.readTimeout, d.writeTimeout)
}
