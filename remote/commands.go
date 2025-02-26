package remote

import (
	"io"
	"net"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	model device.Device

	conn io.ReadWriteCloser
}

// Connect connects to the supplied address.
// The address should be on the format ip:port.
func (c *Control) Connect(address string, model device.Device) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.model = model
	return nil
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
