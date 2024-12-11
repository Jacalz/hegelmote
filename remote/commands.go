package remote

import (
	"io"
	"net"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	conn io.ReadWriteCloser
}

// Connect connects to the supplied address.
// The address should be on the format ip:port.
func (c *Control) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
