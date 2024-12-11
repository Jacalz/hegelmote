package remote

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

var errSorceInputIsZero = errors.New("source indexing starts at 1")

// SetSourceName tells the amplifier to switch to the corresponding source name.
// The input name should match one for the given device type.
func (c *Control) SetSourceName(amp device.Device, name string) error {
	number, err := device.NumberFromName(amp, name)
	if err != nil {
		return err
	}

	return c.SetSourceNumber(number)
}

// SetSourceNumber sets the input source to the given number.
// This will fail if the source number does not exist on the device.
// NOTE: The input source is be indexed from 1.
func (c *Control) SetSourceNumber(number uint) error {
	if number == 0 {
		return errSorceInputIsZero
	}

	parameter := strconv.FormatUint(uint64(number), 10)
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", parameter)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetSourceName returns the currently selected input source.
// The source number will try to map number to a source name on the device type.
func (c *Control) GetSourceName(amp device.Device) (string, error) {
	number, err := c.GetSourceNumber()
	if err != nil {
		return "", err
	}

	return device.NameFromNumber(amp, uint(number))
}

// GetSourceNumber returns the currently selected source number.
func (c *Control) GetSourceNumber() (uint, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", "?")
	if err != nil {
		return 0, err
	}

	resp := [6]byte{}
	n, err := c.conn.Read(resp[:])
	if err != nil {
		return 0, err
	}

	input := resp[3 : n-1]
	number, err := strconv.ParseUint(string(input), 10, 8)
	if err != nil {
		return 0, err
	}

	return uint(number), nil
}
