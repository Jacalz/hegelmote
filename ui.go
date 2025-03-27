package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type remoteUI struct {
	amplifier statefulController
	window    fyne.Window

	// Widgets:
	powerToggle                      *widget.Button
	volumeDisplay                    *widget.Label
	volumeSlider                     *widget.Slider
	volumeMute, volumeDown, volumeUp *widget.Button
	inputSelector                    *widget.Select
}

func (r *remoteUI) refreshPower() {
	if r.amplifier.status.poweredOn {
		r.powerToggle.SetText("Power off")
	} else {
		r.powerToggle.SetText("Power on")
	}
}

func (r *remoteUI) refreshVolumeSlider() {
	r.volumeSlider.OnChangeEnded = nil

	r.volumeSlider.Value = float64(r.amplifier.status.volume)
	r.volumeSlider.OnChanged(r.volumeSlider.Value)

	if r.amplifier.status.poweredOn && !r.amplifier.status.muted {
		enableAndRefresh(r.volumeSlider)
	} else {
		disableAndRefresh(r.volumeSlider)
	}

	r.volumeSlider.OnChangeEnded = r.onVolumeDragEnd
}

func (r *remoteUI) refreshVolumeButtons() {
	if r.amplifier.status.poweredOn {
		r.volumeMute.Enable()
		r.volumeDown.Enable()
		r.volumeUp.Enable()
	} else {
		r.volumeMute.Disable()
		r.volumeDown.Disable()
		r.volumeUp.Disable()
	}
}

func (r *remoteUI) refreshInput() {
	r.inputSelector.OnChanged = nil
	r.inputSelector.Selected = r.amplifier.status.input

	if r.amplifier.status.poweredOn {
		enableAndRefresh(r.inputSelector)
	} else {
		disableAndRefresh(r.inputSelector)
	}

	r.inputSelector.OnChanged = r.onInputSelect
}

func (r *remoteUI) fullRefresh() {
	r.refreshPower()
	r.refreshVolumeSlider()
	r.refreshVolumeButtons()
	r.refreshInput()
}

func (r *remoteUI) onPowerToggle() {
	r.amplifier.togglePower()
	r.fullRefresh()
}

func (r *remoteUI) onVolumeDrag(percentage float64) {
	r.volumeDisplay.SetText(strconv.Itoa(int(percentage)) + "%")
}

func (r *remoteUI) onVolumeDragEnd(percentage float64) {
	r.amplifier.setVolume(uint8(percentage))
	r.refreshVolumeSlider()
}

func (r *remoteUI) onMute() {
	r.amplifier.toggleMute()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onVolumeDown() {
	r.amplifier.volumeDown()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onVolumeUp() {
	r.amplifier.volumeUp()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onInputSelect(input string) {
	r.amplifier.setInput(input)
	r.refreshInput()
}

func buildRemoteUI(command *remote.Control, w fyne.Window) (*remoteUI, fyne.CanvasObject) {
	ui := &remoteUI{window: w, amplifier: statefulController{control: command}}
	ui.amplifier.load()

	ui.powerToggle = &widget.Button{Text: "Toggle power", OnTapped: ui.onPowerToggle}

	ui.volumeDisplay = &widget.Label{Text: "0%"}
	ui.volumeSlider = &widget.Slider{Min: 0, Max: 100, Step: 1, OnChanged: ui.onVolumeDrag, OnChangeEnded: ui.onVolumeDragEnd}

	ui.volumeMute = &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: ui.onMute}

	ui.volumeDown = &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: ui.onVolumeDown}
	ui.volumeUp = &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: ui.onVolumeUp}

	inputLabel := &widget.Label{Text: "Select input:", TextStyle: fyne.TextStyle{Bold: true}}

	inputs, _ := device.GetInputNames(ui.amplifier.control.Model) // TODO: Move this to a connection step.
	ui.inputSelector = &widget.Select{Options: inputs, PlaceHolder: "Select an input", OnChanged: ui.onInputSelect}

	ui.amplifier.load()
	ui.fullRefresh()

	ui.amplifier.trackChanges(
		func(refresh refreshed) {
			switch refresh {
			case refreshPower:
				fyne.Do(ui.fullRefresh)
			case refreshVolume, refreshMute:
				fyne.Do(ui.refreshVolumeSlider)
			case refreshInput:
				fyne.Do(ui.refreshInput)
			}
		},
	)

	return ui, container.NewVBox(
		ui.powerToggle,
		widget.NewSeparator(),
		container.NewVBox(
			container.NewBorder(nil, nil, nil, ui.volumeDisplay, ui.volumeSlider),
			container.NewGridWithColumns(3, ui.volumeMute, ui.volumeDown, ui.volumeUp),
		),
		widget.NewSeparator(),
		inputLabel,
		ui.inputSelector,
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
