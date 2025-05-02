package ui

import (
	_ "embed"
	"strconv"

	"fyne.io/fyne/v2"
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
	amplifier *remote.ControlWithListener
	host      string
	window    fyne.Window

	poweredOn bool
	volume    remote.Volume
	muted     bool
	input     device.Input

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
	if !m.poweredOn {
		text = "Power on"
	}

	m.powerToggle.SetText(text)
}

func (m *mainUI) refreshVolumeSlider() {
	m.volumeSlider.OnChangeEnded = nil

	m.volumeSlider.Value = float64(m.volume)
	m.volumeSlider.OnChanged(m.volumeSlider.Value)

	if m.poweredOn && !m.muted {
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
	setEnabled(m.volumeMute, m.poweredOn)
	setEnabled(m.volumeDown, m.poweredOn)
	setEnabled(m.volumeUp, m.poweredOn)
}

func (m *mainUI) refreshInput() {
	m.inputSelector.OnChanged = nil
	m.inputSelector.Selected = m.inputSelector.Options[m.input-1]

	if m.poweredOn {
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
	on, err := m.amplifier.TogglePower()
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.poweredOn = on
		m.fullRefresh()
	}
}

func (m *mainUI) onVolumeDrag(percentage float64) {
	m.volumeDisplay.SetText(strconv.FormatUint(uint64(percentage), 10) + "%")
}

func (m *mainUI) onVolumeDragEnd(percentage float64) {
	volume, err := m.amplifier.SetVolume(remote.Volume(percentage))
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onMute() {
	muted, err := m.amplifier.ToggleMute()
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.muted = muted
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onVolumeDown() {
	volume, err := m.amplifier.VolumeDown()
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onVolumeUp() {
	volume, err := m.amplifier.VolumeUp()
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.volume = volume
		m.refreshVolumeSlider()
	}
}

func (m *mainUI) onInputSelect(selected string) {
	input, err := m.amplifier.SetInput(device.Input(m.inputSelector.SelectedIndex() + 1)) // #nosec
	showErrorIfNotNil(err, m.window)
	if err == nil {
		m.input = input
		m.refreshInput()
	}
}

func (m *mainUI) onPowerChanged(poweredOn bool) {
	m.poweredOn = poweredOn
	m.fullRefresh()
}

func (m *mainUI) onVolumeChanged(volume remote.Volume) {
	m.volume = volume
	m.refreshVolumeSlider()
}

func (m *mainUI) onMuteChanged(muted bool) {
	m.muted = muted
	m.refreshVolumeSlider()
}

func (m *mainUI) onInputChanged(input device.Input) {
	m.input = input
	m.refreshInput()
}

func (m *mainUI) onError(err error) {
	fyne.LogError("Received error from state tracker", err)
	dialog.ShowError(err, m.window)
}

func (m *mainUI) load() error {
	on, err := m.amplifier.GetPower()
	if err != nil {
		return err
	}

	m.poweredOn = on

	volume, err := m.amplifier.GetVolume()
	if err != nil {
		return err
	}

	m.volume = volume

	muted, err := m.amplifier.GetVolumeMute()
	if err != nil {
		return err
	}

	m.muted = muted

	input, err := m.amplifier.GetInput()
	if err != nil {
		return err
	}

	m.input = input
	return nil
}

// Build sets up and builds the main user interface.
func Build(a fyne.App, w fyne.Window) (*mainUI, fyne.CanvasObject) {
	ui := &mainUI{window: w}
	ui.amplifier = remote.NewControlWithListener(
		ui.onPowerChanged,
		ui.onVolumeChanged,
		ui.onMuteChanged,
		ui.onInputChanged,
		ui.Disconnect,
		ui.onError,
	)

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
