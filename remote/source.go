package remote

import (
	"fmt"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

// SetSourceInput tells the amplifier to switch to the corresponding device input.
func (c *Control) SetSourceInput(amp device.Device, input string) error {
	number := strconv.Itoa(device.InputNumber(amp, input))
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", number)
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

	percentage := resp[3:n]
	// TODO: Convert number to input string from device.

	return string(percentage), nil
}
