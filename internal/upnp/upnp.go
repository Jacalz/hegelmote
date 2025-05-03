package upnp

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/Jacalz/hegelmote/device"
	upnp "github.com/supersonic-app/go-upnpcast/device"
)

// DiscoveredDevice specifies a discovered Hegel amplifier on the network.
type DiscoveredDevice struct {
	Host  string
	Model device.Type
}

// LookUpDevices searches the local network for discoverable devices.
func LookUpDevices() ([]DiscoveredDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	unfiltered, err := upnp.SearchMediaRenderers(ctx, 1)
	if err != nil {
		return nil, err
	}

	devices := []DiscoveredDevice{}
	for _, found := range unfiltered {
		if !strings.HasPrefix(found.FriendlyName, "Hegel") {
			continue
		}

		rawURL, err := url.Parse(found.URL)
		if err != nil {
			continue
		}

		model, err := device.FromString(found.ModelName)
		if err != nil {
			continue
		}

		devices = append(devices, DiscoveredDevice{
			Host:  rawURL.Hostname(),
			Model: model,
		})
	}

	return devices, nil
}
