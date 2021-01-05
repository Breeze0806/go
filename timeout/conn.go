package timeout

import (
	"errors"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewConn(c net.Conn, readTimeout time.Duration, writeTimeout time.Duration) (*Conn, error) {
	if c == nil {
		return nil, errors.New("net conn is nil")
	}

	return &Conn{
		Conn:         c,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}, nil
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.readTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	defer func() {
		if c.readTimeout != 0 {
			c.SetReadDeadline(time.Time{})
		}
	}()
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.writeTimeout != 0 {
		c.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}
	defer func() {
		if c.writeTimeout != 0 {
			c.SetWriteDeadline(time.Time{})
		}
	}()
	return c.Conn.Write(b)
}
