package remote

import (
	"fmt"
	"net"
	"strconv"

	"github.com/Jacalz/hegelmote/device"
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

// SetSourceInput tells the amplifier to switch to the corresponding device input.
func (c *Control) SetSourceInput(amp device.Device, input string) error {
	number := strconv.Itoa(device.InputNumber(amp, input))
	_, err := fmt.Fprintf(c.conn, commandFormat, "i", number)
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
