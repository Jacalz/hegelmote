package ui

import (
	"fmt"
	"image/color"
	"net/netip"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/internal/upnp"
)

func (m *mainUI) connect(host string, model device.Type) error {
	err := m.amplifier.Connect(host, model)
	if err != nil {
		fyne.LogError("Failed to connect to amplifier", err)
		return err
	}

	inputs, err := device.GetInputNames(model)
	if err != nil {
		fyne.LogError("Failed to get input names for model", err)
		return err
	}

	err = m.load()
	if err != nil {
		fyne.LogError("Failed to load initial state", err)
		return err
	}

	m.inputSelector.Options = inputs
	m.host = host
	m.connectionLabel.SetText("Connected")
	m.powerToggle.Enable()
	m.fullRefresh()
	return nil
}

func (m *mainUI) Disconnect() {
	m.powerToggle.Disable()
	err := m.amplifier.Disconnect()
	if err != nil {
		fyne.LogError("Error on disconnecting", err)
	}
	m.connectionLabel.SetText("Disconnected")
}

func (m *mainUI) setUpConnection(prefs fyne.Preferences, w fyne.Window) {
	host := prefs.String("host")
	modelID := prefs.IntWithFallback("model", -1)
	if host != "" && modelID >= 0 && modelID <= int(device.H590) {
		err := m.connect(host, device.Type(modelID)) // #nosec - Range is checked above!
		if err == nil {
			return
		}

		fyne.LogError("Failed to connect to saved connection", err)
		prefs.RemoveValue("host")
		prefs.RemoveValue("model")
	}

	showConnectionDialog(m, w)
}

func handleConnection(host string, model device.Type, remember bool, ui *mainUI) error {
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
	models := &widget.Select{PlaceHolder: "Device type", Options: device.SupportedTypeNames()}
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

func selectFromOneDevice(remote upnp.DiscoveredDevice, ui *mainUI, w fyne.Window) {
	msg := widget.NewRichTextFromMarkdown(fmt.Sprintf("Found **Hegel %s** at **%s**.", device.SupportedTypeNames()[remote.Model], remote.Host))
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, w)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			err := handleConnection(remote.Host, remote.Model, remember.Checked, ui)
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

func selectFromMultipleDevices(remotes []upnp.DiscoveredDevice, ui *mainUI, w fyne.Window) {
	options := make([]string, 0, len(remotes))
	for _, remote := range remotes {
		options = append(options, fmt.Sprintf("Hegel %s \u2013 %s", device.SupportedTypeNames()[remote.Model], remote.Host))
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
			err := handleConnection(remote.Host, remote.Model, remember.Checked, ui)
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
		devices, err := upnp.LookUpDevices()
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
