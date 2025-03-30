package device

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

var supportedDeviceNames = [...]string{"Röst", "H95", "H120", "H190", "H390", "H590"}

// SupportedDeviceNames returns a slice of all supported devices.
func SupportedDeviceNames() []string {
	return supportedDeviceNames[:]
}

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
