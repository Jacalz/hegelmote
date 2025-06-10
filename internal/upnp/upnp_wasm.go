//go:build wasm

package upnp

import (
	"cmp"
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type upnpResponse struct {
	Devices []DiscoveredDevice `json:"devices"`
	Err     error              `json:"error"`
}

// LookUpDevices searches the local network for discoverable devices.
func LookUpDevices() ([]DiscoveredDevice, error) {
	ws, _, err := websocket.Dial(context.Background(), fmt.Sprintf("ws://%s/upnp", Proxy), nil)
	if err != nil {
		return nil, err
	}
	defer ws.CloseNow()

	response := upnpResponse{}
	err = wsjson.Read(context.Background(), ws, &response)
	return response.Devices, cmp.Or(err, response.Err)
}
