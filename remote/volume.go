package remote

import "fmt"

// Volume specifies a volume in the range 0 to 100.
type Volume = uint8

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(volume Volume) (Volume, error) {
	if volume > 100 {
		return 0, fmt.Errorf("invalid volume: %d", volume)
	}

	packet := createNumericalPacket('v', volume)
	return c.sendWithNumericalResponse(packet)
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() (Volume, error) {
	return c.sendWithNumericalResponse([]byte("-v.u\r"))
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() (Volume, error) {
	return c.sendWithNumericalResponse([]byte("-v.d\r"))
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (Volume, error) {
	return c.sendWithNumericalResponse([]byte("-v.?\r"))
}

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) (bool, error) {
	packet := createBooleanPacket('m', mute)
	return c.sendWithBoolResponse(packet)
}

// ToggleVolumeMute toggles the muting of volume.
func (c *Control) ToggleVolumeMute() (bool, error) {
	return c.sendWithBoolResponse([]byte("-m.t\r"))
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	return c.sendWithBoolResponse([]byte("-m.?\r"))
}
