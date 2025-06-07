package remote

// Delay specifies the status of the connection reset.
type Delay struct {
	Minutes Minutes
	Stopped bool
}

// Minutes defines a number of 0 to 255 minutes.
type Minutes = uint8

// SetResetDelay sets a timeout, in minutes, for when to reset.
func (c *Control) SetResetDelay(delay Minutes) (Delay, error) {
	packet := createNumericalPacket('r', delay)
	return c.reset(packet)
}

// StopResetDelay stops the delayed reset from happening.
func (c *Control) StopResetDelay() (Delay, error) {
	return c.reset([]byte("-r.~\r"))
}

// GetResetDelay returns the current delay until reset.
// Returns the delay or a bool indicating if it is stopped or not.
func (c *Control) GetResetDelay() (Delay, error) {
	return c.reset([]byte("-r.?\r"))
}

func (c *Control) reset(packet []byte) (Delay, error) {
	_, err := c.conn.Write(packet)
	if err != nil {
		return Delay{}, err
	}

	buf, err := c.read('r')
	if err != nil {
		return Delay{}, err
	}

	if buf[3] == '~' {
		return Delay{Stopped: true}, nil
	}

	minutes, err := parseUint8FromBuf(buf)
	return Delay{Minutes: minutes}, err
}
