package device

import "slices"

// Type specifies the Hegel amplifier device type to target.
type Type int

const (
	Röst Type = iota
	H95
	H120
	H190
	H390
	H590
	H190V
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
// -1 is returned if the device is not supported.
func FromString(device string) Type {
	return Type(slices.Index(supportedDeviceNames[:], device)) // #nosec
}
