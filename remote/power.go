package remote

// SetPower sets the power to either on or off.
func (c *Control) SetPower(on bool) (bool, error) {
	packet := []byte("-p.0\r")
	if on {
		packet[3] = '1'
	}

	_, err := c.Conn.Write(packet)
	if err != nil {
		return false, err
	}

	return c.parsePowerResponse()
}

// TogglePower toggles the power on and off.
func (c *Control) TogglePower() (bool, error) {
	_, err := c.Conn.Write([]byte("-p.t\r"))
	if err != nil {
		return false, err
	}

	return c.parsePowerResponse()
}

// GetPower returns the current power status.
func (c *Control) GetPower() (bool, error) {
	_, err := c.Conn.Write([]byte("-p.?\r"))
	if err != nil {
		return false, err
	}

	return c.parsePowerResponse()
}

func (c *Control) parsePowerResponse() (bool, error) {
	buf := [len("-p.?\r")]byte{}
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

	if buf[1] != 'p' {
		return false, errUnexpectedResponse
	}

	return buf[3] == '1', nil
}
