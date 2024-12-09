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
func (c *Control) GetResetDelay() (uint8, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", "?")
	return 0, err
}
