//go:build wasm

package remote

import (
	"context"
	"fmt"

	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/internal/upnp"
	"github.com/coder/websocket"
)

func (c *Control) connect(host string, model device.Type) error {
	ws, _, err := websocket.Dial(context.Background(), fmt.Sprintf("ws://%s/proxy", upnp.Proxy), nil)
	if err != nil {
		return err
	}

	err = ws.Write(context.Background(), websocket.MessageText, []byte(host))
	if err != nil {
		return err
	}

	c.conn = &wsWrapper{ws: ws}
	c.deviceType = model
	return nil
}

type wsWrapper struct {
	ws *websocket.Conn
}

func (w *wsWrapper) Read(p []byte) (int, error) {
	_, buf, err := w.ws.Read(context.Background())
	copy(p, buf)
	return len(buf), err
}

func (w *wsWrapper) Write(p []byte) (int, error) {
	return len(p), w.ws.Write(context.Background(), websocket.MessageText, p)
}

func (w *wsWrapper) Close() error {
	return w.ws.Close(websocket.StatusNormalClosure, "")
}
