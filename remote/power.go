package remote

import "fmt"

const powerCommand = "-p.%s\r"

func (c *Control) SetPower(on bool) error {
	state := "0"
	if on {
		state = "1"
	}

	_, err := fmt.Fprintf(c.conn, powerCommand, state)
	return err
}

func (c *Control) TogglePower() error {
	_, err := fmt.Fprintf(c.conn, powerCommand, "t")
	return err
}

func (c *Control) GetPower() (bool, error) {
	_, err := fmt.Fprintf(c.conn, powerCommand, "?")
	if err != nil {
		return false, err
	}

	buf := [4]byte{}
	_, err = c.conn.Read(buf[:])
	if err != nil {
		return false, err
	}

	return buf[3] == '1', nil
}
