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

	_, err := c.Conn.Write(packet)
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

// StopResetDelay stops the delayed reset from happening.
func (c *Control) StopResetDelay() (Delay, error) {
	_, err := c.Conn.Write([]byte("-r.~\r"))
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

// GetResetDelay returns the current delay until reset.
// Returns the delay or a bool indicating if it is stopped or not.
func (c *Control) GetResetDelay() (Delay, error) {
	_, err := c.Conn.Write([]byte("-r.?\r"))
	if err != nil {
		return Delay{}, err
	}

	return c.parseResetResponse()
}

func (c *Control) parseResetResponse() (Delay, error) {
	buf := [len("-r.255\r")]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return Delay{}, err
	}

	if n < 5 {
		return Delay{}, errUnexpectedResponse
	}

	if buf[1] == 'e' {
		return Delay{}, errorFromCode(buf[3])
	}

	if buf[1] != 'r' {
		return Delay{}, errUnexpectedResponse
	}

	if buf[3] == '~' {
		return Delay{Stopped: true}, nil
	}

	delay := buf[3 : n-1]
	number, err := strconv.ParseUint(string(delay), 10, 8)
	return Delay{Minutes: uint8(number)}, err
}
