// Package remote implements remote IP control of supported Hegel amplifiers.
package remote

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

// Control implements remote IP control of supported Hegel amplifiers.
type Control struct {
	deviceType device.Type

	conn io.ReadWriteCloser
}

// Connect connects to the supplied host address. A port should not be specified.
func (c *Control) Connect(host string, model device.Type) error {
	return c.connect(host, model)
}

// Disconnect closes the remote connection.
func (c *Control) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	return err
}

// GetDeviceType returns the device type of the current connection.
func (c *Control) GetDeviceType() device.Type {
	return c.deviceType
}

func (c *Control) read(expectedCommand byte) ([]byte, error) {
	buf := [len("-v.100\r")]byte{}

	n, err := c.conn.Read(buf[:])
	if err != nil {
		return nil, err
	}

	resp, err := c.verifyResponse(buf, n)
	if err != nil {
		return nil, err
	}

	if resp[1] != expectedCommand {
		return nil, fmt.Errorf("unexpected response: %q", buf[:])
	}

	return resp, nil
}

func (c *Control) verifyResponse(buf [7]byte, n int) ([]byte, error) {
	if n < 5 {
		return nil, fmt.Errorf("unexpected response: %q", buf[:])
	}

	if buf[1] == 'e' {
		return nil, errorFromCode(buf[3])
	}

	return buf[:n], nil
}

func (c *Control) sendWithBoolResponse(packet []byte) (bool, error) {
	_, err := c.conn.Write(packet)
	if err != nil {
		return false, err
	}

	return c.parseOnOffValue(packet[1])
}

func (c *Control) parseOnOffValue(command byte) (bool, error) {
	buf, err := c.read(command)
	if err != nil {
		return false, err
	}

	return buf[3] == '1', nil
}

func (c *Control) sendWithNumericalResponse(packet []byte) (uint8, error) {
	_, err := c.conn.Write(packet)
	if err != nil {
		return 0, err
	}

	return c.parseNumberFromResponse(packet[1])
}

func (c *Control) parseNumberFromResponse(command byte) (uint8, error) {
	buf, err := c.read(command)
	if err != nil {
		return 0, err
	}

	return parseUint8FromBuf(buf)
}

func parseUint8FromBuf(buf []byte) (uint8, error) {
	number := uint16(0)
	for i := 3; i < len(buf)-1 && buf[i] != '\r'; i++ {
		char := buf[i]
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid uint8 value: %q", buf)
		}

		number = number*10 + uint16(char-'0')
	}

	if number > 255 {
		return 0, fmt.Errorf("value %d does not fit in uint8", number)
	}

	return uint8(number), nil
}

func createNumericalPacket(command byte, value uint8) []byte {
	packet := make([]byte, 0, 7)
	packet = append(packet, '-', command, '.')
	packet = strconv.AppendUint(packet, uint64(value), 10)
	return append(packet, '\r')
}

func createBooleanPacket(command byte, value bool) []byte {
	packet := []byte{'-', command, '.', '0', '\r'}
	if value {
		packet[3] = '1'
	}
	return packet
}

// Mapping of error values. Index zero corresponds to error 1 and so on.
// The following error codes were reverse engineered by sending incorrect commands.
func errorFromCode(code byte) error {
	switch code - '0' {
	case 1:
		return errors.New("malformed command")
	case 2:
		return errors.New("unknown command")
	case 3:
		return errors.New("invalid parameter")
	default:
		return fmt.Errorf("unexpected error code: %d", code)
	}
}
