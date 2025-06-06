package remote

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) (bool, error) {
	packet := []byte("-m.0\r")
	if mute {
		packet[3] = '1'
	}

	return c.mute(packet)
}

// ToggleVolumeMute toggles the muting of volume.
func (c *Control) ToggleVolumeMute() (bool, error) {
	return c.mute([]byte("-m.t\r"))
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	return c.mute([]byte("-m.?\r"))
}

func (c *Control) mute(packet []byte) (bool, error) {
	_, err := c.conn.Write(packet)
	if err != nil {
		return false, err
	}

	return c.parseOnOffValue('m')
}
