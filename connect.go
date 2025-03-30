package main

import (
	"context"
	"fmt"
	"image/color"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	upnp "github.com/supersonic-app/go-upnpcast/device"
)

type discoveredDevice struct {
	host      string
	model     device.Device
	modelName string
}

func lookUpDevices() ([]discoveredDevice, error) {
	unfiltered, err := upnp.SearchMediaRenderers(context.TODO(), 1)
	if err != nil {
		return nil, err
	}

	devices := []discoveredDevice{}
	for _, found := range unfiltered {
		_, ok := strings.CutPrefix(found.FriendlyName, "Hegel")
		if !ok {
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
			host:      rawURL.Hostname(),
			modelName: found.ModelName,
			model:     model,
		})
	}

	return devices, nil
}

func selectFromOneDevice(remote discoveredDevice, ui *remoteUI, w fyne.Window) {
	msg := widget.NewRichTextFromMarkdown(fmt.Sprintf("Found **Hegel %s** at **%s**.", remote.modelName, remote.host))
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, w)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			if remember.Checked {
				prefs := fyne.CurrentApp().Preferences()
				prefs.SetString("host", remote.host)
				prefs.SetInt("model", int(remote.model))
			}

			ui.connect(remote.host, remote.model)
			connectionDialog.Hide()
		},
	}
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func selectFromMultipleDevices(remotes []discoveredDevice, ui *remoteUI, w fyne.Window) {
	options := make([]string, 0, len(remotes))
	for _, remote := range remotes {
		options = append(options, fmt.Sprintf("Hegel %s \u2013 %s", remote.modelName, remote.host))
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
			if remember.Checked {
				prefs := fyne.CurrentApp().Preferences()
				prefs.SetString("host", remote.host)
				prefs.SetInt("model", int(remote.model))
			}

			ui.connect(remote.host, remote.model)
			connectionDialog.Hide()
		},
	}
	connect.Disable()
	selection.OnChanged = func(_ string) { connect.Enable() }
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func showConnectionDialog(ui *remoteUI, w fyne.Window) {
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
			dialog.ShowInformation("No devices found", "The LAN did not seem to contain any supported devices.", w)
			return
		}

		if len(devices) == 1 {
			selectFromMultipleDevices(devices, ui, w)
			return
		}

		selectFromOneDevice(devices[0], ui, w)
	}()
}
