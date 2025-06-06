//go:build !wasm

package remote

import (
	"context"
	"net"
	"time"

	"github.com/Jacalz/hegelmote/device"
)

func (c *Control) connect(host string, model device.Type) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", host+":50001")
	if err != nil {
		return err
	}

	c.conn = conn
	c.deviceType = model
	return nil
}
