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
		return 0, fmt.Errorf("invalid volume value: %d", volume)
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-v."...)
	packet = strconv.AppendUint(packet, uint64(volume), 10)
	packet = append(packet, '\r')

	_, err := c.conn.Write(packet)
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() (Volume, error) {
	_, err := c.conn.Write([]byte("-v.u\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() (Volume, error) {
	_, err := c.conn.Write([]byte("-v.d\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (Volume, error) {
	_, err := c.conn.Write([]byte("-v.?\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

func (c *Control) parseVolumeResponse() (Volume, error) {
	return c.parseNumberFromResponse('v')
}
