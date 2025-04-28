package remote

import (
	"errors"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

var errSorceInputIsZero = errors.New("source indexing starts at 1")

// SetSourceName tells the amplifier to switch to the corresponding source name.
// The input name should match one for the given device type.
func (c *Control) SetSourceName(name string) error {
	number, err := device.InputFromName(c.Model, name)
	if err != nil {
		return err
	}

	return c.SetSourceNumber(number)
}

// SetSourceNumber sets the input source to the given number.
// This will fail if the source number does not exist on the device.
func (c *Control) SetSourceNumber(number device.Input) error {
	if number == 0 {
		return errSorceInputIsZero
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-i."...)
	packet = strconv.AppendUint(packet, uint64(number), 10)
	packet = append(packet, '\r')

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetSourceName returns the currently selected input source.
// The source number will try to map number to a source name on the device type.
func (c *Control) GetSourceName() (string, error) {
	input, err := c.GetSourceNumber()
	if err != nil {
		return "", err
	}

	return device.NameFromNumber(c.Model, input)
}

// GetSourceNumber returns the currently selected source number.
func (c *Control) GetSourceNumber() (device.Input, error) {
	_, err := c.Conn.Write([]byte("-i.?\r"))
	if err != nil {
		return 0, err
	}

	buf := [6]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return 0, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return 0, err
	}

	input := buf[3 : n-1]
	number, err := strconv.ParseUint(string(input), 10, 8)
	return device.Input(number), err
}
