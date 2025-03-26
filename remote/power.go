package remote

// SetPower sets the power to either on or off.
func (c *Control) SetPower(on bool) error {
	packet := []byte("-p.0\r")
	if on {
		packet[3] = '1'
	}

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// TogglePower toggles the power on and off.
func (c *Control) TogglePower() error {
	_, err := c.Conn.Write([]byte("-p.t\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetPower returns the current power status.
func (c *Control) GetPower() (bool, error) {
	_, err := c.Conn.Write([]byte("-p.?\r"))
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

	return buf[3] == '1', nil
}
