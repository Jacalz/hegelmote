package remote

import (
	"errors"
	"fmt"
	"strconv"
)

var errInvalidPercentage = errors.New("invalid percentage value")

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(percentage uint) error {
	if percentage > 100 {
		return errInvalidPercentage
	}

	value := strconv.FormatUint(uint64(percentage), 10)
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", value)
	return err
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "u")
	return err
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "d")
	return err
}

func (c *Control) GetVolume() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", "?")
	return err
}

// SetVolumeMute allows turning on or off mute.
func (c *Control) SetVolumeMute(mute bool) error {
	state := "0"
	if mute {
		state = "1"
	}

	_, err := fmt.Fprintf(c.conn, commandFormat, "m", state)
	return err
}

func (c *Control) GetVolumeMute() error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "m", "?")
	return err
}
