package device

import (
	"errors"
	"slices"
)

// Device defines the Hegel amplifier to target.
type Device = uint

const (
	Röst Device = iota
	H95
	H120
	H190
	H390
	H590
)

var deviceInputs = [...][]string{InputsRöst, InputsH95, InputsH120, InputsH190, InputsH390, InputsH590}

var (
	errInvalidDevice = errors.New("invalid device type")
	errInvalidInput  = errors.New("input not on device")
)

// FromString takes a model name as string and returns the corresponding [Device] ID for it.
func FromString(device string) (Device, error) {
	switch device {
	case "Röst":
		return Röst, nil
	case "H95":
		return H95, nil
	case "H120":
		return H120, nil
	case "H190":
		return H190, nil
	case "H390":
		return H390, nil
	case "H590":
		return H590, nil
	}

	return Device(^uint(0)), errInvalidDevice
}

// GetInputs returns the list of input names for the given device.
func GetInputNames(device Device) ([]string, error) {
	if device > H590 {
		return nil, errInvalidDevice
	}

	return deviceInputs[device], nil
}

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

	return uint(number) + 1, nil // #nosec: Known input!
}

// NameFromNumber returns the corresponding input name for the input number.
// NOTE: The input is indexed from 1.
func NameFromNumber(device Device, number uint) (string, error) {
	if device > H590 {
		return "", errInvalidDevice
	}

	inputs := deviceInputs[device]
	if number > uint(len(inputs)) {
		return "", errInvalidInput
	}

	return inputs[number-1], nil
}
