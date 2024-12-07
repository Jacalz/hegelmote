package remote

import "fmt"

func (c *Control) VolumeControl(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", param)
	return err
}

func (c *Control) VolumeMute(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "m", param)
	return err
}
