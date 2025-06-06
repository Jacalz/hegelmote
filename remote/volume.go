package remote

import (
	"fmt"
	"strconv"
)

// Volume specifies a volume in the range 0 to 100.
type Volume = uint8

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(volume Volume) (Volume, error) {
	if volume > 100 {
		return 0, fmt.Errorf("invalid volume: %d", volume)
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-v."...)
	packet = strconv.AppendUint(packet, uint64(volume), 10)
	packet = append(packet, '\r')

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
