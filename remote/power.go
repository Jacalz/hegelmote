package remote

// SetPower sets the power to either on or off.
func (c *Control) SetPower(on bool) (bool, error) {
	packet := []byte("-p.0\r")
	if on {
		packet[3] = '1'
	}

	return c.power(packet)
}

// TogglePower toggles the power on and off.
func (c *Control) TogglePower() (bool, error) {
	return c.power([]byte("-p.t\r"))
}

// GetPower returns the current power status.
func (c *Control) GetPower() (bool, error) {
	return c.power([]byte("-p.?\r"))
}

func (c *Control) power(packet []byte) (bool, error) {
	_, err := c.conn.Write(packet)
	if err != nil {
		return false, err
	}

	return c.parseOnOffValue('p')
}
