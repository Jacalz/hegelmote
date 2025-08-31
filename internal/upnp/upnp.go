// Package upnp provides tooling to discover Hegel amplifiers on the network.
package upnp

import "github.com/Jacalz/hegelmote/device"

// Proxy specifies the proxy server to use for UPnP discovery.
var Proxy = "localhost:8086"

// DiscoveredDevice specifies a discovered Hegel amplifier on the network.
type DiscoveredDevice struct {
	Host  string
	Model device.Type
}
