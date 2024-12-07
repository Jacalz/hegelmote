package remote

import "fmt"

func (c *Control) ResetConnection(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", param)
	return err
}
