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
