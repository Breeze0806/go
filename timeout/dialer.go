package timeout

import (
	"context"
	"net"
	"time"
)

//普通拨号器
type dialer interface {
	//根据连接类型network和地址address生成连接
	dial(network, address string) (net.Conn, error)
	//根据上下文ctx，连接类型network和地址address生成连接
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

//Dialer 超时网络连接拨号器
type Dialer struct {
	dialer                     // 普通网络连接拨号器
	readTimeout  time.Duration //读超时
	writeTimeout time.Duration //写超时
}

//NewDialer 根据读超时readTimeout以及写writeTimeout超时网络生成连接拨号器
func NewDialer(readTimeout time.Duration, writeTimeout time.Duration) *Dialer {
	return &Dialer{
		dialer:       &defaultDialer{},
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
	}
}

//Dial 根据连接类型network和地址address生成超时连接
//如果读超时和写超时为0，则生成普通连接
//连接类型可以是tcp，udp等等，address支持ipv4,ipv8或者ip:port等等
func (d *Dialer) Dial(network, address string) (conn net.Conn, err error) {
	conn, err = d.dial(network, address)
	if err != nil {
		return nil, err
	}
	return d.conn(conn)
}

//DialContext 根据上下文ctx，连接类型network和地址address生成超时连接
//如果读超时和写超时为0，则生成普通连接
//连接类型可以是tcp，udp等等，address支持ipv4,ipv8或者ip:port等等
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
