package device

import (
	"errors"
	"slices"
)

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

var (
	errInvalidDevice = errors.New("invalid device type")
	errInvalidInput  = errors.New("input not on device")
)

// NumberFromName returns the corresponding input number for the input name.
// NOTE: The output is indexed from 1.
func NumberFromName(device Device, input string) (uint, error) {
	if device > H590 {
		return 0, errInvalidDevice
	}

	inputs := deviceInputs[device]
	number := slices.Index(inputs, input)
	if number == -1 {
		return 0, errInvalidInput
	}

	return uint(number) + 1, nil
}

// NameFromNumber returns the corresponding input name for the input number.
// NOTE: The input is indexed from 1.
func NameFromNumber(device Device, number uint) (string, error) {
	if device > H590 {
		return "", errInvalidDevice
	}

	inputs := deviceInputs[device]
	if number > uint(len(deviceInputs)) {
		return "", errInvalidInput
	}

	return inputs[number-1], nil
}
