package remote

import (
	"net"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	Model device.Device

	Conn net.Conn
}

// Connect connects to the supplied host address. A port should not be specified.
func (c *Control) Connect(host string, model device.Device) error {
	conn, err := net.Dial("tcp", host+":50001")
	if err != nil {
		return err
	}

	c.Conn = conn
	c.Model = model
	return nil
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.Conn == nil {
		return nil
	}

	return c.Conn.Close()
}

// Read provides access to reading out the raw data from the command buffer.
// This can be used to listen for changes sent by the amplifier.
func (c *Control) Read() ([]byte, error) {
	buf := [len("-v.100\r")]byte{}

	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}
