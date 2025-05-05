package device

import "slices"

var supportedDeviceNames = [...]string{"Röst", "H95", "H120", "H190", "H390", "H590"}

// Type specifies the Hegel amplifier device type to target.
type Type int

const (
	Röst Type = iota
	H95
	H120
	H190
	H390
	H590
)

// String returns the string name of the device.
func (t Type) String() string {
	if !IsSupported(t) {
		return ""
	}

	return supportedDeviceNames[t]
}

// IsSupported reports true if the given device type is supported.
func IsSupported(device Type) bool {
	return device >= Röst && device < Type(len(supportedDeviceNames))
}

// SupportedTypeNames returns a slice of all supported devices.
func SupportedTypeNames() []string {
	return supportedDeviceNames[:]
}

// FromString takes a model name as string and returns the corresponding [Type] ID for it.
func FromString(device string) (Type, error) {
	index := slices.Index(supportedDeviceNames[:], device)
	if index == -1 {
		return -1, errInvalidDevice
	}

	return Type(index), nil // #nosec
}
