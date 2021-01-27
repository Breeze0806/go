package timeout

import (
	"errors"
	"net"
	"time"
)

//Conn 超时连接
type Conn struct {
	net.Conn                   //普通连接
	readTimeout  time.Duration //读超时
	writeTimeout time.Duration //写超时
}

//NewConn 通过普通连接c, 读超时readTimeout，写超时writeTimeout生成查实连接
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

//Read 读取网络字节流b，设置了对应的读超时，返回读取字符个数n和读取错误err
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

//Write 写入网络字节流b，设置了对应的写超时，返回写入字符个数n和写入错误err
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
