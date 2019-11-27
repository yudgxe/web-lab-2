package server

import (
	"net"
	"time"
)

type conn struct {
	net.Conn
	IdleTimeout time.Duration
}

func (c *conn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Write(p)
}

func (c *conn) Read(b []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Read(b)
}

func (c *conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
