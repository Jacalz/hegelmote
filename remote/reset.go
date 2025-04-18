package remote

import (
	"strconv"
)

// SetResetDelay sets a timeout, in minutes, for when to reset.
func (c *Control) SetResetDelay(delay uint8) error {
	packet := make([]byte, 0, 7)
	packet = append(packet, "-r."...)
	packet = strconv.AppendUint(packet, uint64(delay), 10)
	packet = append(packet, '\r')

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// StopResetDelay stops the delayed reset from happening.
func (c *Control) StopResetDelay() error {
	_, err := c.Conn.Write([]byte("-r.~\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetResetDelay returns the current delay until reset.
// Returns the delay or a bool indicating if it is stopped or not.
func (c *Control) GetResetDelay() (uint8, bool, error) {
	_, err := c.Conn.Write([]byte("-r.?\r"))
	if err != nil {
		return 0, false, err
	}

	buf := [7]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return 0, false, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return 0, false, err
	}

	// Check if reset is stopped or not enabled.
	if n >= 4 && buf[3] == '~' {
		return 0, true, nil
	}

	number := buf[3 : n-1]
	delay, err := strconv.ParseUint(string(number), 10, 8)
	return uint8(delay), false, err
}
