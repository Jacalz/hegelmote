//go:build wasm

package remote

import (
	"context"

	"github.com/Jacalz/hegelmote/device"
	"github.com/coder/websocket"
)

func (c *Control) connect(host string, model device.Type) error {
	ws, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/proxy", nil)
	if err != nil {
		return err
	}

	err = ws.Write(context.Background(), websocket.MessageText, []byte(host))
	if err != nil {
		return err
	}

	c.conn = websocket.NetConn(context.Background(), ws, websocket.MessageText)
	c.deviceType = model
	return nil
}
