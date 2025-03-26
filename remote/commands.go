package remote

import (
	"net"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	model device.Device

	Conn net.Conn
}

// Connect connects to the supplied address.
// The address should be on the format ip:port.
func (c *Control) Connect(address string, model device.Device) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.Conn = conn
	c.model = model
	return nil
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.Conn == nil {
		return nil
	}

	return c.Conn.Close()
}

func (c *Control) Read() ([]byte, error) {
	buf := [len("-v.100\r")]byte{}

	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
