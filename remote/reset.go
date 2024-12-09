package remote

import (
	"fmt"
	"strconv"
)

// SetResetDelay sets a timeout, in minutes, for when to reset.
func (c *Control) SetResetDelay(delay uint8) error {
	number := strconv.FormatInt(int64(delay), 10)
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", number)
	return err
}

// StopResetDelay stops the delayed reset from happening.
func (c *Control) StopResetDelay() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", "~")
	return err
}

// GetResetDelay returns the current delay until reset.
// Returns the delay or a bool indicating if it is stopped or not.
func (c *Control) GetResetDelay() (uint8, bool, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", "?")

	resp := [6]byte{}
	n, err := c.conn.Read(resp[:])
	if err != nil {
		return 0, false, err
	}

	// Check if reset is stopped or not enabled.
	if n == 4 && resp[3] == '~' {
		return 0, true, nil
	}

	number := resp[3:n]
	delay, err := strconv.ParseUint(string(number), 10, 8)
	return uint8(delay), false, err
}
