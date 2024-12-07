package remote

import (
	"fmt"
	"net"
)

const commandFormat = "-%s.%s\r"

type Control struct {
	Ip   string
	Port string

	conn net.Conn
}

func (c *Control) Connect(ip, port string) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Control) Disconnect() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *Control) SourceInput(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", param)
	return err
}

func (c *Control) VolumeControl(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "v", param)
	return err
}

func (c *Control) VolumeMute(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "m", param)
	return err
}

func (c *Control) ResetConnection(param string) error {
	_, err := fmt.Fprintf(c.conn, commandFormat, "r", param)
	return err
}
