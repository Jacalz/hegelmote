package remote

// SetPower sets the power to either on or off.
func (c *Control) SetPower(on bool) (bool, error) {
	packet := createBooleanPacket('p', on)
	return c.sendWithBoolResponse(packet)
}

// TogglePower toggles the power on and off.
func (c *Control) TogglePower() (bool, error) {
	return c.sendWithBoolResponse([]byte("-p.t\r"))
}

// GetPower returns the current power status.
func (c *Control) GetPower() (bool, error) {
	return c.sendWithBoolResponse([]byte("-p.?\r"))
}
