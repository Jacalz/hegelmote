package device

import (
	"errors"
	"slices"
)

// Input is a number specifying which input to use.
// Note that this is indexed from one and not zero.
type Input = uint8

var (
	errInvalidDevice = errors.New("invalid device type")
	errInvalidInput  = errors.New("unsupported input for device")
)

// GetInputs returns the list of input names for the given device.
func GetInputNames(device Type) ([]string, error) {
	switch device {
	case Röst:
		return InputsRöst[:], nil
	case H95:
		return InputsH95[:], nil
	case H120:
		return InputsH120[:], nil
	case H190:
		return InputsH190[:], nil
	case H390:
		return InputsH390[:], nil
	case H590:
		return InputsH590[:], nil
	}

	return nil, errInvalidDevice
}

// InputFromName returns the corresponding input number for the input name.
// NOTE: The output is indexed from 1.
func InputFromName(device Type, input string) (Input, error) {
	inputs, err := GetInputNames(device)
	if err != nil {
		return 0, err
	}

	number := slices.Index(inputs, input)
	if number == -1 {
		return 0, errInvalidInput
	}

	return Input(number) + 1, nil // #nosec
}

// NameFromNumber returns the corresponding input name for the input number.
func NameFromNumber(device Type, input Input) (string, error) {
	inputs, err := GetInputNames(device)
	if err != nil {
		return "", err
	}

	if int(input) > len(inputs) {
		return "", errInvalidInput
	}

	return inputs[input-1], nil
}
