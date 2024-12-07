package remote

import (
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
