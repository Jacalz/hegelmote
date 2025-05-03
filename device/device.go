package device

import "slices"

// Type specifies the Hegel amplifier device type to target.
type Type = uint

const (
	Röst Type = iota
	H95
	H120
	H190
	H390
	H590
	unsupported
)

var supportedDeviceNames = [...]string{"Röst", "H95", "H120", "H190", "H390", "H590"}

// SupportedTypeNames returns a slice of all supported devices.
func SupportedTypeNames() []string {
	return supportedDeviceNames[:]
}

// FromString takes a model name as string and returns the corresponding [Type] ID for it.
func FromString(device string) (Type, error) {
	index := slices.Index(supportedDeviceNames[:], device)
	if index == -1 {
		return unsupported, errInvalidDevice
	}

	return Type(index), nil // #nosec
}
