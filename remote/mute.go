package remote

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) (bool, error) {
	packet := []byte("-m.0\r")
	if mute {
		packet[3] = '1'
	}

	_, err := c.Conn.Write(packet)
	if err != nil {
		return false, err
	}

	return c.parseMuteResponse()
}

// ToggleVolumeMute toggles the muting of volume.
func (c *Control) ToggleVolumeMute() (bool, error) {
	_, err := c.Conn.Write([]byte("-m.t\r"))
	if err != nil {
		return false, err
	}

	return c.parseMuteResponse()
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	_, err := c.Conn.Write([]byte("-m.?\r"))
	if err != nil {
		return false, err
	}

	return c.parseMuteResponse()
}

func (c *Control) parseMuteResponse() (bool, error) {
	buf := [len("-m.?\r")]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return false, err
	}

	if n < 5 {
		return false, errUnexpectedResponse
	}

	if buf[1] == 'e' {
		return false, errorFromCode(buf[3])
	}

	if buf[1] != 'm' {
		return false, errUnexpectedResponse
	}

	return buf[3] == '1', nil
}
