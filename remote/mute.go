package remote

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) (bool, error) {
	packet := createBooleanPacket('m', mute)
	return c.sendWithBoolResponse(packet)
}

// ToggleVolumeMute toggles the muting of volume.
func (c *Control) ToggleVolumeMute() (bool, error) {
	return c.sendWithBoolResponse([]byte("-m.t\r"))
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	return c.sendWithBoolResponse([]byte("-m.?\r"))
}
