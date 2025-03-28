package remote

import (
	"errors"
	"strconv"
)

var errInvalidPercentage = errors.New("invalid percentage value")

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(percentage uint8) error {
	if percentage > 100 {
		return errInvalidPercentage
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-v."...)
	packet = strconv.AppendUint(packet, uint64(percentage), 10)
	packet = append(packet, '\r')

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() error {
	_, err := c.Conn.Write([]byte("-v.u\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() error {
	_, err := c.Conn.Write([]byte("-v.d\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (uint, error) {
	_, err := c.Conn.Write([]byte("-v.?\r"))
	if err != nil {
		return 0, err
	}

	buf := [7]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return 0, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return 0, err
	}

	volume := buf[3 : n-1]
	percentage, err := strconv.ParseUint(string(volume), 10, 8)
	return uint(percentage), err
}

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) error {
	packet := []byte("-m.0\r")
	if mute {
		packet[3] = '1'
	}

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// ToggleVolumeMute toggles the muting of volume.
func (c *Control) ToggleVolumeMute() error {
	_, err := c.Conn.Write([]byte("-m.t\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	_, err := c.Conn.Write([]byte("-m.?\r"))
	if err != nil {
		return false, err
	}

	buf := [5]byte{}
	_, err = c.Conn.Read(buf[:])
	if err != nil {
		return false, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return false, err
	}

	return buf[3] == '1', err
}
