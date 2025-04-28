package device

import (
	"errors"
	"slices"
)

// Input is a number specifying which input to use.
// Note that this is indexed from one and not zero.
type Input = uint8

var deviceInputs = [...][]string{InputsRöst, InputsH95, InputsH120, InputsH190, InputsH390, InputsH590}

var (
	errInvalidDevice = errors.New("invalid device type")
	errInvalidInput  = errors.New("input not on device")
)

// GetInputs returns the list of input names for the given device.
func GetInputNames(device Device) ([]string, error) {
	if device > H590 {
		return nil, errInvalidDevice
	}

	return deviceInputs[device], nil
}

// InputFromName returns the corresponding input number for the input name.
// NOTE: The output is indexed from 1.
func InputFromName(device Device, input string) (Input, error) {
	if device > H590 {
		return 0, errInvalidDevice
	}

	inputs := deviceInputs[device]
	number := slices.Index(inputs, input)
	if number == -1 {
		return 0, errInvalidInput
	}

	return Input(number) + 1, nil // #nosec: Known input!
}

// NameFromNumber returns the corresponding input name for the input number.
func NameFromNumber(device Device, number Input) (string, error) {
	if device > H590 {
		return "", errInvalidDevice
	}

	inputs := deviceInputs[device]
	if int(number) > len(inputs) {
		return "", errInvalidInput
	}

	return inputs[number-1], nil
}
