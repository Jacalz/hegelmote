package ui

import (
	"context"
	"fmt"
	"image/color"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	upnp "github.com/supersonic-app/go-upnpcast/device"
)

type discoveredDevice struct {
	host  string
	model device.Device
}

func lookUpDevices() ([]discoveredDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	unfiltered, err := upnp.SearchMediaRenderers(ctx, 1)
	if err != nil {
		return nil, err
	}

	devices := []discoveredDevice{}
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

		devices = append(devices, discoveredDevice{
			host:  rawURL.Hostname(),
			model: model,
		})
	}

	return devices, nil
}

func handleConnection(host string, model device.Device, remember bool, ui *mainUI) error {
	err := ui.connect(host, model)
	if err != nil {
		fyne.LogError("Failed to connect", err)
		return err
	}

	if remember && model <= device.H590 {
		prefs := fyne.CurrentApp().Preferences()
		prefs.SetString("host", host)
		prefs.SetInt("model", int(model)) // #nosec - Checked by model <= device.H590 above!
	}
	return nil
}

func selectManually(ui *mainUI, w fyne.Window) {
	hostname := &widget.Entry{PlaceHolder: "IP Address (no port)"}
	models := &widget.Select{PlaceHolder: "Device type", Options: device.SupportedDeviceNames()}
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(hostname, models, remember)

	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, w)
	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			model, _ := device.FromString(models.Selected)
			err := handleConnection(hostname.Text, model, remember.Checked, ui)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			connectionDialog.Hide()
		},
	}

	connect.Disable()
	hostname.OnChanged = func(_ string) {
		_, errIP := netip.ParseAddr(hostname.Text)
		hasModel := models.SelectedIndex() != -1
		if errIP == nil && hasModel {
			connect.Enable()
			return
		}

		connect.Disable()
	}
	models.OnChanged = hostname.OnChanged

	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func selectFromOneDevice(remote discoveredDevice, ui *mainUI, w fyne.Window) {
	msg := widget.NewRichTextFromMarkdown(fmt.Sprintf("Found **Hegel %s** at **%s**.", device.SupportedDeviceNames()[remote.model], remote.host))
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, w)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			err := handleConnection(remote.host, remote.model, remember.Checked, ui)
			if err != nil {
				selectManually(ui, w)
				return
			}
			connectionDialog.Hide()
		},
	}
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func selectFromMultipleDevices(remotes []discoveredDevice, ui *mainUI, w fyne.Window) {
	options := make([]string, 0, len(remotes))
	for _, remote := range remotes {
		options = append(options, fmt.Sprintf("Hegel %s \u2013 %s", device.SupportedDeviceNames()[remote.model], remote.host))
	}

	msg := &widget.Label{Text: "Multiple devices were discovered:"}
	selection := &widget.Select{PlaceHolder: "Choose a device", Options: options}
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, selection, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, w)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			index := selection.SelectedIndex()
			if index == -1 {
				return
			}

			remote := remotes[index]
			err := handleConnection(remote.host, remote.model, remember.Checked, ui)
			if err != nil {
				selectManually(ui, w)
				return
			}
			connectionDialog.Hide()
		},
	}
	connect.Disable()
	selection.OnChanged = func(_ string) { connect.Enable() }
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func showConnectionDialog(ui *mainUI, w fyne.Window) {
	prop := canvas.NewRectangle(color.Transparent)
	prop.SetMinSize(fyne.NewSquareSize(75))

	activity := widget.NewActivity()
	activity.Start()
	d := dialog.NewCustomWithoutButtons("Looking for amplifiers on LAN\u2026", container.NewStack(prop, activity), w)
	d.SetOnClosed(activity.Stop)
	d.Show()

	go func() {
		defer d.Hide()
		devices, err := lookUpDevices()
		if err != nil || len(devices) == 0 {
			fyne.LogError("Failed to search for devices", err)
			selectManually(ui, w)
			return
		}

		if len(devices) > 1 {
			selectFromMultipleDevices(devices, ui, w)
			return
		}

		selectFromOneDevice(devices[0], ui, w)
	}()
}
