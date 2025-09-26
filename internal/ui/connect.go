package ui

import (
	"fmt"
	"image/color"
	"net/netip"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/internal/upnp"
)

func (m *mainUI) connect(host string, model device.Type) error {
	err := m.amplifier.Connect(host, model)
	if err != nil {
		return err
	}

	inputs, err := device.GetInputNames(model)
	if err != nil {
		return err
	}

	err = m.load()
	if err != nil {
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

func (m *mainUI) setUpConnection() {
	prefs := fyne.CurrentApp().Preferences()
	host := prefs.String("host")
	model := device.Type(prefs.IntWithFallback("model", -1))
	if host == "" || !device.IsSupported(model) {
		m.showConnectionDialog()
		return
	}

	go func() {
		err := m.connect(host, model)
		if err != nil {
			fyne.LogError("Failed to connect to remembered connection", err)
			fyne.Do(func() {
				prefs.RemoveValue("host")
				prefs.RemoveValue("model")
				m.showConnectionDialog()
			})
		}
	}()
}

func (m *mainUI) handleConnection(host string, model device.Type, remember bool) error {
	err := m.connect(host, model)
	if err != nil {
		fyne.LogError("Failed to connect", err)
		return err
	}

	if remember && device.IsSupported(model) {
		prefs := fyne.CurrentApp().Preferences()
		prefs.SetString("host", host)
		prefs.SetInt("model", int(model))
	}
	return nil
}

func (m *mainUI) showManualConnectionDialog() {
	hostname := &widget.Entry{PlaceHolder: "IP Address (no port)"}
	models := &widget.Select{PlaceHolder: "Device type", Options: device.SupportedTypeNames()}
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(hostname, models, remember)

	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, m.window)
	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			model := device.Type(models.SelectedIndex()) // #nosec
			err := m.handleConnection(hostname.Text, model, remember.Checked)
			if err != nil {
				dialog.ShowError(err, m.window)
				return
			}

			connectionDialog.Hide()
		},
	}

	connect.Disable()
	hostname.OnChanged = func(_ string) {
		_, errIP := netip.ParseAddr(hostname.Text)
		hasDeviceType := models.SelectedIndex() != -1
		setEnabled(connect, errIP == nil && hasDeviceType)
	}
	models.OnChanged = hostname.OnChanged

	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func (m *mainUI) showConnectOneDialog(remote upnp.DiscoveredDevice) {
	msg := widget.NewRichTextFromMarkdown(fmt.Sprintf("Found **Hegel %s** at **%s**.", remote.Model.String(), remote.Host))
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, m.window)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			err := m.handleConnection(remote.Host, remote.Model, remember.Checked)
			if err != nil {
				m.showManualConnectionDialog()
			}
			connectionDialog.Hide()
		},
	}
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func (m *mainUI) showConnectMultipleDialog(remotes []upnp.DiscoveredDevice) {
	options := make([]string, 0, len(remotes))
	for _, remote := range remotes {
		options = append(options, fmt.Sprintf("Hegel %s \u2013 %s", remote.Model.String(), remote.Host))
	}

	msg := &widget.Label{Text: "Multiple devices were discovered:"}
	selection := &widget.Select{PlaceHolder: "Choose a device", Options: options}
	remember := &widget.Check{Text: "Remember connection"}
	content := container.NewVBox(msg, selection, remember)
	connectionDialog := dialog.NewCustomWithoutButtons("Connect to device", content, m.window)

	connect := &widget.Button{
		Text:       "Connect",
		Importance: widget.HighImportance,
		OnTapped: func() {
			index := selection.SelectedIndex()
			if index == -1 {
				return
			}

			remote := remotes[index]
			err := m.handleConnection(remote.Host, remote.Model, remember.Checked)
			if err != nil {
				m.showManualConnectionDialog()
			}
			connectionDialog.Hide()
		},
	}
	connect.Disable()
	selection.OnChanged = func(_ string) { connect.Enable() }
	connectionDialog.SetButtons([]fyne.CanvasObject{connect})
	fyne.Do(connectionDialog.Show)
}

func (m *mainUI) showConnectionDialog() {
	prop := canvas.NewRectangle(color.Transparent)
	prop.SetMinSize(fyne.NewSquareSize(75))

	activity := widget.NewActivity()
	d := dialog.NewCustomWithoutButtons("Looking for amplifiers on LAN\u2026", container.NewStack(prop, activity), m.window)
	d.SetOnClosed(activity.Stop)

	go func() {
		if runtime.GOOS == "js" {
			m.selectProxyServer()
		}

		devices, err := upnp.LookUpDevices()
		if err != nil {
			fyne.LogError("Failed to search for devices", err)
			devices = nil // Zero length so we show manual connection dialog.
		}

		d.Hide()
		switch len(devices) {
		case 0:
			m.showManualConnectionDialog()
		case 1:
			m.showConnectOneDialog(devices[0])
		default:
			m.showConnectMultipleDialog(devices)
		}
	}()

	activity.Start()
	d.Show()
}

func (m *mainUI) onConnectionInfo() {
	info := &widget.Form{Items: []*widget.FormItem{
		{Text: "Address", Widget: &widget.Label{Text: m.host}},
		{Text: "Model", Widget: &widget.Label{Text: "Hegel " + m.amplifier.GetDeviceType().String()}},
		{Text: "Status", Widget: &widget.Label{Text: m.connectionLabel.Text}},
	}}

	prefs := fyne.CurrentApp().Preferences()
	var infoDialog *dialog.CustomDialog

	disconnect := &widget.Button{Text: "Disconnect", Icon: theme.CancelIcon(), Importance: widget.LowImportance, OnTapped: func() {
		infoDialog.Hide()
		m.Disconnect()
		m.showConnectionDialog()
	}}

	forget := &widget.Button{Text: "Forget", Icon: theme.MediaReplayIcon(), Importance: widget.LowImportance}
	forget.OnTapped = func() {
		forget.Disable()
		prefs.RemoveValue("host")
		prefs.RemoveValue("model")
	}

	host := prefs.String("host")
	model := device.Type(prefs.IntWithFallback("model", -1))
	if host == "" || !device.IsSupported(model) {
		forget.Disable()
	}

	prop := &canvas.Rectangle{}
	prop.SetMinSize(fyne.NewSquareSize(theme.Padding()))

	infoDialog = dialog.NewCustom("Connection info", "Dismiss", container.NewVBox(info, container.NewGridWithRows(1, disconnect, forget), prop), m.window)
	infoDialog.Show()
}

func (m *mainUI) selectProxyServer() {
	prefs := fyne.CurrentApp().Preferences()
	proxy := prefs.String("proxy")
	if proxy == "" {
		wait := make(chan struct{})
		proxyEntry := &widget.Entry{PlaceHolder: "localhost:8086"}
		done := &widget.Button{Text: "Confirm", Importance: widget.HighImportance, OnTapped: func() { close(wait) }}
		dialog.ShowCustomWithoutButtons("Select a proxy server to use", container.NewVBox(proxyEntry, &widget.Separator{}, container.NewHBox(done)), m.window)
		<-wait

		proxy = proxyEntry.Text
		prefs.SetString("proxy", proxy)
	}
	upnp.Proxy = proxy
}
