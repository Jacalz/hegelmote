package remote

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
