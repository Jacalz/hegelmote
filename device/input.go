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
func GetInputNames(device Type) ([]string, error) {
	if device >= unsupported {
		return nil, errInvalidDevice
	}

	return deviceInputs[device], nil
}

// InputFromName returns the corresponding input number for the input name.
// NOTE: The output is indexed from 1.
func InputFromName(device Type, input string) (Input, error) {
	if device >= unsupported {
		return 0, errInvalidDevice
	}

	inputs := deviceInputs[device]
	number := slices.Index(inputs, input)
	if number == -1 {
		return 0, errInvalidInput
	}

	return Input(number) + 1, nil // #nosec
}

// NameFromNumber returns the corresponding input name for the input number.
func NameFromNumber(device Type, input Input) (string, error) {
	if device >= unsupported {
		return "", errInvalidDevice
	}

	inputs := deviceInputs[device]
	if int(input) > len(inputs) {
		return "", errInvalidInput
	}

	return inputs[input-1], nil
}
