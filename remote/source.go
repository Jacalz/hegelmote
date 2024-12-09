package remote

import (
	"fmt"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

// SetSourceInput tells the amplifier to switch to the corresponding device input.
func (c *Control) SetSourceInput(amp device.Device, input string) error {
	number, err := device.NumberFromName(amp, input)
	if err != nil {
		return err
	}
	parameter := strconv.FormatUint(uint64(number), 10)
	_, err = fmt.Fprintf(c.conn, commandFormat, "i", parameter)
	return err
}

// GetSourceInput returns the currently selected input source.
func (c *Control) GetSourceInput(amp device.Device) (string, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", "?")
	if err != nil {
		return "", err
	}

	resp := [6]byte{}
	n, err := c.conn.Read(resp[:])
	if err != nil {
		return "", err
	}

	input := resp[3 : n-1]
	number, err := strconv.ParseUint(string(input), 10, 8)
	if err != nil {
		return "", err
	}

	return device.NameFromNumber(amp, uint(number))
}
