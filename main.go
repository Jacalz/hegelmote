package main

import (
	_ "embed"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type remoteUI struct {
	amplifier state
	window    fyne.Window

	// Widgets:
	powerToggle                      *widget.Button
	volumeDisplay                    *widget.Label
	volumeSlider                     *widget.Slider
	volumeMute, volumeDown, volumeUp *widget.Button
	inputSelector                    *widget.Select
}

func (r *remoteUI) syncState() {
	r.volumeSlider.OnChangeEnded = nil
	r.inputSelector.OnChanged = nil

	// Power:
	if r.amplifier.poweredOn {
		r.powerToggle.SetText("Power off")
		r.volumeMute.Enable()
		r.volumeDown.Enable()
		r.volumeUp.Enable()
		r.inputSelector.Enable()
	} else {
		r.powerToggle.SetText("Power on")
		r.volumeMute.Disable()
		r.volumeDown.Disable()
		r.volumeUp.Disable()
		r.inputSelector.Disable()
	}

	// Volume:
	r.volumeSlider.Value = float64(r.amplifier.volume)
	r.volumeSlider.OnChanged(r.volumeSlider.Value)

	// Mute:
	if r.amplifier.muted || !r.amplifier.poweredOn {
		if r.volumeSlider.Disabled() {
			r.volumeSlider.Refresh()
		}
		r.volumeSlider.Disable()
	} else {
		if !r.volumeSlider.Disabled() {
			r.volumeSlider.Refresh()
		}
		r.volumeSlider.Enable()
	}

	// Input:
	r.inputSelector.SetSelected(r.amplifier.input)

	r.inputSelector.OnChanged = r.onInputSelect
	r.volumeSlider.OnChangeEnded = r.onVolumeDragEnd
}

func (r *remoteUI) onPowerToggle() {
	r.amplifier.togglePower()
	r.syncState()
}

func (r *remoteUI) onVolumeDrag(percentage float64) {
	r.volumeDisplay.SetText(strconv.Itoa(int(percentage)) + "%")
}

func (r *remoteUI) onVolumeDragEnd(percentage float64) {
	r.amplifier.setVolume(uint8(percentage))
	r.syncState()
}

func (r *remoteUI) onMute() {
	r.amplifier.toggleMute()
	r.syncState()
}

func (r *remoteUI) onVolumeDown() {
	r.amplifier.volumeDown()
	r.syncState()
}

func (r *remoteUI) onVolumeUp() {
	r.amplifier.volumeUp()
	r.syncState()
}

func (r *remoteUI) onInputSelect(input string) {
	r.amplifier.setInput(input)
	r.syncState()
}

func buildRemoteUI(command *remote.Control, w fyne.Window) (*remoteUI, fyne.CanvasObject) {
	ui := &remoteUI{window: w, amplifier: state{control: command}}
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
	ui.syncState()

	ui.amplifier.listenForChanges(func() { fyne.Do(ui.syncState) })

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

func main() {
	a := app.NewWithID("io.github.jacalz.hegelmote")
	w := a.NewWindow("Hegelmote")

	command := &remote.Control{}

	err := command.Connect("192.168.1.251:50001", device.H95)
	if err != nil {
		panic(err)
	}

	ui, content := buildRemoteUI(command, w)
	defer ui.amplifier.disconnect()

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 400))
	w.SetMaster()
	w.ShowAndRun()
}
