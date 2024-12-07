package device

import "slices"

// Device defines the Hegel amplifier to target.
type Device uint8

const (
	Röst Device = iota
	H95
	H120
	H190
	H390
	H590
)

var deviceInputs = [][]string{InputsRöst, InputsH95, InputsH120, InputsH190, InputsH390, InputsH590}

// InputNumber returns the corresponding input number for the input name.
func InputNumber(device Device, input string) int {
	inputs := deviceInputs[device]
	number := slices.Index(inputs, input)
	if number == -1 {
		return number
	}

	return number + 1
}
