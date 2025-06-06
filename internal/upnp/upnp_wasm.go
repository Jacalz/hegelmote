//go:build wasm

package upnp

import (
	"context"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// LookUpDevices searches the local network for discoverable devices.
func LookUpDevices() ([]DiscoveredDevice, error) {
	ws, _, err := websocket.Dial(context.Background(), "ws://localhost:8080/upnp", nil)
	if err != nil {
		return nil, err
	}
	defer ws.CloseNow()

	devices := []DiscoveredDevice{}
	err = wsjson.Read(context.Background(), ws, &devices)
	return devices, err
}
