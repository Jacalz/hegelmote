package remote

import "strconv"

// Volume specifies a volume in the range 0 to 100.
type Volume = uint8

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(volume Volume) error {
	if volume > 100 {
		return errInvalidVolume
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-v."...)
	packet = strconv.AppendUint(packet, uint64(volume), 10)
	packet = append(packet, '\r')

	_, err := c.Conn.Write(packet)
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() error {
	_, err := c.Conn.Write([]byte("-v.u\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() error {
	_, err := c.Conn.Write([]byte("-v.d\r"))
	if err != nil {
		return err
	}

	return c.parseErrorResponse()
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (Volume, error) {
	_, err := c.Conn.Write([]byte("-v.?\r"))
	if err != nil {
		return 0, err
	}

	buf := [7]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return 0, err
	}

	err = parseErrorFromBuffer(buf[:])
	if err != nil {
		return 0, err
	}

	volume := buf[3 : n-1]
	percentage, err := strconv.ParseUint(string(volume), 10, 8)
	return Volume(percentage), err
}
