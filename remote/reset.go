package remote

import (
	"strconv"
)

// Delay specifies the status of the connection reset.
type Delay struct {
	Minutes Minutes
	Stopped bool
}

// Minutes defines a number of 0 to 255 minutes.
type Minutes = uint8

// SetResetDelay sets a timeout, in minutes, for when to reset.
func (c *Control) SetResetDelay(delay Minutes) (Delay, error) {
	packet := make([]byte, 0, 7)
	packet = append(packet, "-r."...)
	packet = strconv.AppendUint(packet, uint64(delay), 10)
	packet = append(packet, '\r')

	_, err := c.conn.Write(packet)
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

// StopResetDelay stops the delayed reset from happening.
func (c *Control) StopResetDelay() (Delay, error) {
	_, err := c.conn.Write([]byte("-r.~\r"))
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

// GetResetDelay returns the current delay until reset.
// Returns the delay or a bool indicating if it is stopped or not.
func (c *Control) GetResetDelay() (Delay, error) {
	_, err := c.conn.Write([]byte("-r.?\r"))
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

func (c *Control) parseResetResponse() (Delay, error) {
	buf, err := c.readCommand('r')
	if err != nil {
		return Delay{}, err
	}

	if buf[3] == '~' {
		return Delay{Stopped: true}, nil
	}

	minutes, err := parseUint8FromBuf(buf)
	return Delay{Minutes: minutes}, err
}
