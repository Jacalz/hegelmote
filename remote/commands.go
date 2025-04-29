package remote

import (
	"net"
	"strconv"

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

	if n < 5 {
		return nil, errUnexpectedResponse
	}

	if buf[1] == 'e' {
		return nil, errorFromCode(buf[3])
	}

	return buf[:n], nil
}

func (c *Control) read(expectedCommand byte) ([]byte, error) {
	buf, err := c.Read()
	if err != nil {
		return nil, err
	}

	if buf[1] != expectedCommand {
		return nil, errUnexpectedResponse
	}

	return buf, nil
}

func parseUint8FromBuf(buf []byte) (uint8, error) {
	str := buf[3 : len(buf)-1]
	number, err := strconv.ParseUint(string(str), 10, 8)
	return uint8(number), err
}
