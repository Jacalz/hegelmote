package remote

import (
	"errors"
	"fmt"
	"strconv"
)

var errInvalidPercentage = errors.New("invalid percentage value")

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(percentage uint8) error {
	if percentage > 100 {
		return errInvalidPercentage
	}

	value := strconv.FormatUint(uint64(percentage), 10)
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", value)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "u")
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "d")
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (uint, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "?")
	if err != nil {
		return 0, err
	}

	buf := [7]byte{}
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return 0, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return 0, err
	}

	volume := buf[3 : n-1]
	percentage, err := strconv.ParseUint(string(volume), 10, 8)
	return uint(percentage), err
}

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) error {
	state := "0"
	if mute {
		state = "1"
	}

	_, err := fmt.Fprintf(c.conn, commandFormat, "m", state)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetVolumeMute returns true if the device is muted.
func (c *Control) GetVolumeMute() (bool, error) {
	_, err := fmt.Fprintf(c.conn, commandFormat, "m", "?")
	if err != nil {
		return false, err
	}

	buf := [5]byte{}
	_, err = c.conn.Read(buf[:])
	if err != nil {
		return false, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return false, err
	}

	return buf[3] == '1', err
}
