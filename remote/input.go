package remote

import (
	"strconv"

	"github.com/Jacalz/hegelmote/device"
)

// SetInputFromName tells the amplifier to switch to the corresponding source name.
// The input name should match one for the given device type.
func (c *Control) SetInputFromName(name string) (device.Input, error) {
	number, err := device.InputFromName(c.Model, name)
	if err != nil {
		return 0, err
	}

	return c.SetInput(number)
}

// SetInput sets the input source to the given number.
// This will fail if the source number does not exist on the device.
func (c *Control) SetInput(number device.Input) (device.Input, error) {
	if number == 0 {
		return 0, errInputIsZero
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-i."...)
	packet = strconv.AppendUint(packet, uint64(number), 10)
	packet = append(packet, '\r')

	_, err := c.conn.Write(packet)
	if err != nil {
		return 0, err
	}

	return c.parseInputResponse()
}

// GetInputName returns the currently selected input source.
// The source number will try to map number to a source name on the device type.
func (c *Control) GetInputName() (string, error) {
	input, err := c.GetInput()
	if err != nil {
		return "", err
	}

	return device.NameFromNumber(c.Model, input)
}

// GetInput returns the currently selected source number.
func (c *Control) GetInput() (device.Input, error) {
	_, err := c.conn.Write([]byte("-i.?\r"))
	if err != nil {
		return 0, err
	}

	return c.parseInputResponse()
}

func (c *Control) parseInputResponse() (device.Input, error) {
	return c.parseNumberFromResponse('i')
}
