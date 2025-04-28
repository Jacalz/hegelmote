package ui

import (
	_ "embed"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Jacalz/hegelmote/assets/img"
	"github.com/Jacalz/hegelmote/device"
)

type mainUI struct {
	amplifier statefulController
	host      string
	current   state
	window    fyne.Window

	// Widgets:
	powerToggle                      *widget.Button
	volumeLabel, volumeDisplay       *widget.Label
	volumeSlider                     *widget.Slider
	volumeMute, volumeDown, volumeUp *widget.Button
	inputLabel                       *widget.Label
	inputSelector                    *widget.Select
	connectionLabel                  *widget.Label
	connectionInfoButton             *widget.Button
}

func (m *mainUI) refreshPower() {
	if m.current.poweredOn {
		m.powerToggle.SetText("Power off")
	} else {
		m.powerToggle.SetText("Power on")
	}
}

func (m *mainUI) refreshVolumeSlider() {
	m.volumeSlider.OnChangeEnded = nil

	m.volumeSlider.Value = float64(m.current.volume)
	m.volumeSlider.OnChanged(m.volumeSlider.Value)

	if m.current.poweredOn && !m.current.muted {
		setLabelImportance(m.volumeLabel, widget.MediumImportance)
		enableAndRefresh(m.volumeSlider)
		setLabelImportance(m.volumeDisplay, widget.MediumImportance)
	} else {
		setLabelImportance(m.volumeLabel, widget.LowImportance)
		disableAndRefresh(m.volumeSlider)
		setLabelImportance(m.volumeDisplay, widget.LowImportance)
	}

	m.volumeSlider.OnChangeEnded = m.onVolumeDragEnd
}

func (m *mainUI) refreshVolumeButtons() {
	if m.current.poweredOn {
		m.volumeMute.Enable()
		m.volumeDown.Enable()
		m.volumeUp.Enable()
	} else {
		m.volumeMute.Disable()
		m.volumeDown.Disable()
		m.volumeUp.Disable()
	}
}

func (m *mainUI) refreshInput() {
	m.inputSelector.OnChanged = nil
	m.inputSelector.Selected = m.current.input

	if m.current.poweredOn {
		setLabelImportance(m.inputLabel, widget.MediumImportance)
		enableAndRefresh(m.inputSelector)
	} else {
		setLabelImportance(m.inputLabel, widget.LowImportance)
		disableAndRefresh(m.inputSelector)
	}

	m.inputSelector.OnChanged = m.onInputSelect
}

func (m *mainUI) fullRefresh() {
	m.refreshPower()
	m.refreshVolumeSlider()
	m.refreshVolumeButtons()
	m.refreshInput()
}

func (m *mainUI) onPowerToggle() {
	m.current = m.amplifier.togglePower()
	m.fullRefresh()
}

func (m *mainUI) onVolumeDrag(percentage float64) {
	m.volumeDisplay.SetText(strconv.FormatUint(uint64(percentage), 10) + "%")
}

func (m *mainUI) onVolumeDragEnd(percentage float64) {
	m.current = m.amplifier.setVolume(uint8(percentage))
	m.refreshVolumeSlider()
}

func (m *mainUI) onMute() {
	m.current = m.amplifier.toggleMute()
	m.refreshVolumeSlider()
}

func (m *mainUI) onVolumeDown() {
	m.current = m.amplifier.volumeDown()
	m.refreshVolumeSlider()
}

func (m *mainUI) onVolumeUp() {
	m.current = m.amplifier.volumeUp()
	m.refreshVolumeSlider()
}

func (m *mainUI) onInputSelect(input string) {
	m.current = m.amplifier.setInput(input)
	m.refreshInput()
}

func (m *mainUI) onConnectionInfo() {
	info := &widget.Form{Items: []*widget.FormItem{
		{Text: "Address", Widget: &widget.Label{Text: m.host}},
		{Text: "Model", Widget: &widget.Label{Text: "Hegel " + device.SupportedDeviceNames()[m.amplifier.Model]}},
		{Text: "Status", Widget: &widget.Label{Text: m.connectionLabel.Text}},
	}}

	prefs := fyne.CurrentApp().Preferences()
	var infoDialog *dialog.CustomDialog

	disconnect := &widget.Button{Text: "Disconnect", Icon: theme.CancelIcon(), Importance: widget.LowImportance, OnTapped: func() {
		infoDialog.Hide()
		m.Disconnect()
		showConnectionDialog(m, m.window)
		m.amplifier.closing = false
	}}

	forget := &widget.Button{Text: "Forget", Icon: theme.MediaReplayIcon(), Importance: widget.LowImportance}
	forget.OnTapped = func() {
		forget.Disable()
		prefs.RemoveValue("host")
		prefs.RemoveValue("model")
	}

	host := prefs.String("host")
	modelID := prefs.IntWithFallback("model", -1)
	if host == "" || modelID == -1 {
		forget.Disable()
	}

	prop := &canvas.Rectangle{}
	prop.SetMinSize(fyne.NewSquareSize(theme.Padding()))

	infoDialog = dialog.NewCustom("Connection info", "Dismiss", container.NewVBox(info, container.NewGridWithRows(1, disconnect, forget), prop), m.window)
	infoDialog.Show()
}

func (m *mainUI) connect(host string, model device.Device) error {
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

	m.inputSelector.Options = inputs
	m.current = m.amplifier.load()
	m.host = host
	m.connectionLabel.SetText("Connected")
	m.powerToggle.Enable()
	m.fullRefresh()

	m.amplifier.trackChanges(
		func(refresh refreshed, newState state) {
			fyne.Do(func() {
				m.current = newState

				switch refresh {
				case refreshPower:
					m.fullRefresh()
				case refreshVolume, refreshMute:
					m.refreshVolumeSlider()
				case refreshInput:
					m.refreshInput()
				case reset:
					m.Disconnect()
				}
			})
		},
	)

	m.amplifier.runResetLoop()
	return nil
}

func (m *mainUI) Disconnect() {
	m.powerToggle.Disable()
	m.amplifier.disconnect()
	m.connectionLabel.SetText("Disconnected")
}

func (m *mainUI) setUpConnection(prefs fyne.Preferences, w fyne.Window) {
	host := prefs.String("host")
	modelID := prefs.IntWithFallback("model", -1)
	if host != "" && modelID >= 0 && modelID <= int(device.H590) {
		err := m.connect(host, device.Device(modelID)) // #nosec - Range is checked above!
		if err == nil {
			return
		}

		fyne.LogError("Failed to connect to saved connection", err)
		prefs.RemoveValue("host")
		prefs.RemoveValue("model")
	}

	showConnectionDialog(m, w)
}

// Build sets up and builds the main user interface.
func Build(a fyne.App, w fyne.Window) (*mainUI, fyne.CanvasObject) {
	ui := &mainUI{window: w}

	ui.powerToggle = &widget.Button{Icon: img.PowerIcon, Text: "Toggle power", OnTapped: ui.onPowerToggle}

	ui.volumeLabel = &widget.Label{Text: "Change volume:", TextStyle: fyne.TextStyle{Bold: true}}
	ui.volumeDisplay = &widget.Label{Text: "0%"}
	ui.volumeSlider = &widget.Slider{Min: 0, Max: 100, Step: 1, OnChanged: ui.onVolumeDrag, OnChangeEnded: ui.onVolumeDragEnd}
	ui.volumeMute = &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: ui.onMute}
	ui.volumeDown = &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: ui.onVolumeDown}
	ui.volumeUp = &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: ui.onVolumeUp}

	ui.inputLabel = &widget.Label{Text: "Select input:", TextStyle: fyne.TextStyle{Bold: true}}
	ui.inputSelector = &widget.Select{PlaceHolder: "Select an input", OnChanged: ui.onInputSelect}

	ui.connectionLabel = &widget.Label{Text: "Disconnected", Truncation: fyne.TextTruncateEllipsis}
	ui.connectionInfoButton = &widget.Button{Icon: theme.InfoIcon(), Importance: widget.LowImportance, OnTapped: ui.onConnectionInfo}

	ui.setUpConnection(a.Preferences(), w)

	return ui, container.NewVBox(
		ui.powerToggle,
		widget.NewSeparator(),
		ui.volumeLabel,
		container.NewBorder(nil, nil, nil, ui.volumeDisplay, ui.volumeSlider),
		container.NewGridWithColumns(3, ui.volumeMute, ui.volumeDown, ui.volumeUp),
		widget.NewSeparator(),
		ui.inputLabel,
		ui.inputSelector,
		layout.NewSpacer(),
		container.NewBorder(nil, nil, nil, ui.connectionInfoButton, ui.connectionLabel),
	)
}

type disableableWidget interface {
	fyne.Widget
	fyne.Disableable
}

func enableAndRefresh(wid disableableWidget) {
	if !wid.Disabled() {
		wid.Refresh()
	}
	wid.Enable()
}

func disableAndRefresh(wid disableableWidget) {
	if wid.Disabled() {
		wid.Refresh()
	}
	wid.Disable()
}

func setLabelImportance(label *widget.Label, importance widget.Importance) {
	if label.Importance == importance {
		return
	}

	label.Importance = importance
	label.Refresh()
}
