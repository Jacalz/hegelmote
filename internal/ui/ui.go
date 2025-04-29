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
	"github.com/Jacalz/hegelmote/remote"
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
	text := "Power off"
	if !m.current.poweredOn {
		text = "Power on"
	}

	m.powerToggle.SetText(text)
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
	setEnabled(m.volumeMute, m.current.poweredOn)
	setEnabled(m.volumeDown, m.current.poweredOn)
	setEnabled(m.volumeUp, m.current.poweredOn)
}

func (m *mainUI) refreshInput() {
	m.inputSelector.OnChanged = nil
	m.inputSelector.Selected = m.inputSelector.Options[m.current.input-1]

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
	on, err := m.amplifier.togglePower()
	if err == nil {
		m.current.poweredOn = on
		m.fullRefresh()
	}
}

func (m *mainUI) onVolumeDrag(percentage float64) {
	m.volumeDisplay.SetText(strconv.FormatUint(uint64(percentage), 10) + "%")
}

func (m *mainUI) onVolumeDragEnd(percentage float64) {
	volume, err := m.amplifier.setVolume(remote.Volume(percentage))
	if err == nil {
		m.current.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onMute() {
	muted, err := m.amplifier.toggleMute()
	if err == nil {
		m.current.muted = muted
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onVolumeDown() {
	volume, err := m.amplifier.volumeDown()
	if err == nil {
		m.current.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onVolumeUp() {
	volume, err := m.amplifier.volumeUp()
	if err == nil {
		m.current.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onInputSelect(selected string) {
	input, err := m.amplifier.setInput(device.Input(m.inputSelector.SelectedIndex() + 1))
	if err == nil {
		m.current.input = input
		m.refreshInput()
	}
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

func (m *mainUI) load() error {
	on, err := m.amplifier.GetPower()
	if err != nil {
		return err
	}

	m.current.poweredOn = on

	volume, err := m.amplifier.GetVolume()
	if err != nil {
		return err
	}

	m.current.volume = volume

	muted, err := m.amplifier.GetVolumeMute()
	if err != nil {
		return err
	}

	m.current.muted = muted

	input, err := m.amplifier.GetInput()
	if err != nil {
		return err
	}

	m.current.input = input
	return nil
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

	m.amplifier.trackChanges(
		func(refresh refreshed, newState state) {
			fyne.Do(func() {
				switch refresh {
				case refreshPower:
					m.current.poweredOn = newState.poweredOn
					m.fullRefresh()
				case refreshVolume:
					m.current.volume = newState.volume
					m.refreshVolumeSlider()
				case refreshMute:
					m.current.muted = newState.muted
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
