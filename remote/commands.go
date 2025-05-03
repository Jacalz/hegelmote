package remote

import (
	"net"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	deviceType device.Type

	conn net.Conn
}

// Connect connects to the supplied host address. A port should not be specified.
func (c *Control) Connect(host string, model device.Type) error {
	conn, err := net.Dial("tcp", host+":50001")
	if err != nil {
		return err
	}

	c.conn = conn
	c.deviceType = model
	return nil
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// GetDeviceType returns the device type of the current connection.
func (c *Control) GetDeviceType() device.Type {
	return c.deviceType
}

func (c *Control) read() ([]byte, error) {
	buf := [len("-v.100\r")]byte{}

	n, err := c.conn.Read(buf[:])
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

func (c *Control) readCommand(expectedCommand byte) ([]byte, error) {
	buf, err := c.read()
	if err != nil {
		return nil, err
	}

	if buf[1] != expectedCommand {
		return nil, errUnexpectedResponse
	}

	return buf, nil
}

func (c *Control) parseOnOffValue(command byte) (bool, error) {
	buf, err := c.readCommand(command)
	if err != nil {
		return false, err
	}

	return buf[3] == '1', nil
}

func (c *Control) parseNumberFromResponse(command byte) (uint8, error) {
	buf, err := c.readCommand(command)
	if err != nil {
		return 0, err
	}

	return parseUint8FromBuf(buf)
}

func parseUint8FromBuf(buf []byte) (uint8, error) {
	str := string(buf[3 : len(buf)-1])
	number, err := strconv.ParseUint(str, 10, 8)
	return uint8(number), err
}
