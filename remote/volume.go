package remote

import "strconv"

// Volume specifies a volume in the range 0 to 100.
type Volume = uint8

// SetVolume sets the volume to a value between 0 and 100.
func (c *Control) SetVolume(volume Volume) (Volume, error) {
	if volume > 100 {
		return 0, errInvalidVolume
	}

	packet := make([]byte, 0, 7)
	packet = append(packet, "-v."...)
	packet = strconv.AppendUint(packet, uint64(volume), 10)
	packet = append(packet, '\r')

	_, err := c.Conn.Write(packet)
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// VolumeUp increases the volume one step.
func (c *Control) VolumeUp() (Volume, error) {
	_, err := c.Conn.Write([]byte("-v.u\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// VolumeDown decreases the volume one step.
func (c *Control) VolumeDown() (Volume, error) {
	_, err := c.Conn.Write([]byte("-v.d\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

// GetVolume returns the currrently selected volume percentage.
func (c *Control) GetVolume() (Volume, error) {
	_, err := c.Conn.Write([]byte("-v.?\r"))
	if err != nil {
		return 0, err
	}

	return c.parseVolumeResponse()
}

func (c *Control) parseVolumeResponse() (Volume, error) {
	buf := [len("-m.100\r")]byte{}
	n, err := c.Conn.Read(buf[:])
	if err != nil {
		return 0, err
	}

	if n < 5 {
		return 0, errUnexpectedResponse
	}

	if buf[1] == 'e' {
		return 0, errorFromCode(buf[3])
	}

	if buf[1] != 'v' {
		return 0, errUnexpectedResponse
	}

	volume := buf[3 : n-1]
	percentage, err := strconv.ParseUint(string(volume), 10, 8)
	return Volume(percentage), err
}
