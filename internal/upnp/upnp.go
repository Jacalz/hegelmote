package upnp

import "github.com/Jacalz/hegelmote/device"

// DiscoveredDevice specifies a discovered Hegel amplifier on the network.
type DiscoveredDevice struct {
	Host  string
	Model device.Type
}
