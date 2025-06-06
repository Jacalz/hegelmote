//go:build wasm

package remote

import (
	"context"
	"time"

	"github.com/Jacalz/hegelmote/device"
	"github.com/coder/websocket"
)

func (c *Control) connect(host string, model device.Type) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/proxy", nil)
	if err != nil {
		return err
	}

	c.conn = websocket.NetConn(context.Background(), conn, websocket.MessageText)
	c.deviceType = model
	return nil
}
