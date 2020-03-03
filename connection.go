package main

import (
	"github.com/juju/ratelimit"
	"net"
)

type RateLimitedConn struct {
	net.Conn
	*ratelimit.Bucket
}

func (wrap RateLimitedConn) Read(b []byte) (n int, err error) {
	written, err := wrap.Conn.Read(b)

	wrap.Bucket.Wait(int64(written))

	return written, err
}

func (wrap RateLimitedConn) Write(b []byte) (n int, err error) {
	wrap.Bucket.Wait(int64(len(b)))

	return wrap.Conn.Write(b)
}
