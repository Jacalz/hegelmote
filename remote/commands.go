package remote

import (
	"net"
	"time"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	model device.Device

	conn net.Conn
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

func (c *Control) Read() ([]byte, error) {
	buf := [len("-v.100\r")]byte{}

	c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	defer c.conn.SetReadDeadline(time.Time{})

	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
