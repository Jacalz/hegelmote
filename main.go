package main

import (
	_ "embed"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type remoteUI struct {
	// State:
	control *remote.Control
	model   device.Device
	window  fyne.Window

	// Widgets:
	powerToggle *widget.Button

	volumeDisplay                    *widget.Label
	volumeSlider                     *widget.Slider
	volumeMute, volumeDown, volumeUp *widget.Button

	inputSelector *widget.Select
}

func (r *remoteUI) setupSync() {
	go r.runBackgroundSync()
}

func (r *remoteUI) runBackgroundSync() {
	ticker := time.NewTicker(500 * time.Millisecond)

	sync := func() {
		current := readState(r.control, r.model)
		if current == nil {
			ticker.Stop()
			return
		}

		fyne.Do(func() {
			r.syncState(current)
		})
	}

	sync()

	for range ticker.C {
		sync()
	}
}

func (r *remoteUI) syncState(current *state) {
	// Power:
	if current.poweredOn {
		r.powerToggle.SetText("Power off")
	} else {
		r.powerToggle.SetText("Power on")
	}

	// Mute:
	if current.muted {
		r.volumeSlider.Disable()
	} else {
		r.volumeSlider.Enable()
	}

	// Volume:
	r.volumeSlider.SetValue(float64(current.volume))

	// Input:
	r.inputSelector.Options = current.inputs
	r.inputSelector.SetSelected(current.input)
}

func (r *remoteUI) onPowerToggle() {
	err := r.control.TogglePower()
	if err != nil {
		fyne.LogError("Failed to toggle power", err)
	}
}

func (r *remoteUI) onVolumeDrag(percentage float64) {
	r.volumeDisplay.SetText(strconv.Itoa(int(percentage)) + "%")
}

func (r *remoteUI) onVolumeDragEnd(percentage float64) {
	err := r.control.SetVolume(uint8(percentage))
	if err != nil {
		fyne.LogError("Failed to set volume", err)
	}
}

func (r *remoteUI) onMute() {
	err := r.control.ToggleVolumeMute()
	if err != nil {
		fyne.LogError("Failed to toggle mute", err)
	}
}

func (r *remoteUI) onVolumeDown() {
	err := r.control.VolumeDown()
	if err != nil {
		fyne.LogError("Failed to lower volume", err)
	}
}

func (r *remoteUI) onVolumeUp() {
	err := r.control.VolumeUp()
	if err != nil {
		fyne.LogError("Failed to increase volume", err)
	}
}

func (r *remoteUI) onInputSelect(input string) {
	err := r.control.SetSourceName(r.model, input)
	if err != nil {
		fyne.LogError("Failed to set input", err)
	}
}

func buildRemoteUI(command *remote.Control, w fyne.Window) fyne.CanvasObject {
	ui := remoteUI{window: w, control: command, model: device.H95}
	defer ui.setupSync()

	ui.powerToggle = &widget.Button{Text: "Toggle power", OnTapped: ui.onPowerToggle}

	ui.volumeDisplay = &widget.Label{Text: "0%"}
	ui.volumeSlider = &widget.Slider{Min: 0, Max: 100, Step: 1, OnChanged: ui.onVolumeDrag, OnChangeEnded: ui.onVolumeDragEnd}

	ui.volumeMute = &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: ui.onMute}

	ui.volumeDown = &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: ui.onVolumeDown}
	ui.volumeUp = &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: ui.onVolumeUp}

	inputLabel := &widget.Label{Text: "Select input:", TextStyle: fyne.TextStyle{Bold: true}}

	ui.inputSelector = &widget.Select{PlaceHolder: "Select an input", OnChanged: ui.onInputSelect}

	return container.NewVBox(
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
	defer command.Disconnect()

	err := command.Connect("192.168.1.251:50001")
	if err != nil {
		panic(err)
	}

	w.SetContent(buildRemoteUI(command, w))
	w.Resize(fyne.NewSize(300, 400))
	w.SetMaster()
	w.ShowAndRun()
}
