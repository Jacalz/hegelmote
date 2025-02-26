package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

func showErrorOnFailure(err error, w fyne.Window) {
	if err != nil {
		dialog.ShowError(err, w)
	}
}

func buildRemoteUI(command *remote.Control, w fyne.Window) fyne.CanvasObject {
	power := &widget.Button{}
	power.OnTapped = func() {
		showErrorOnFailure(command.TogglePower(), w)

		on, err := command.GetPower()
		showErrorOnFailure(err, w)
		if on {
			power.SetText("Power off")
		} else {
			power.SetText("Power on")
		}
	}

	on, err := command.GetPower()
	showErrorOnFailure(err, w)
	if on {
		power.Text = "Power off"
	} else {
		power.Text = "Power on"
	}

	volume, err := command.GetVolume()
	showErrorOnFailure(err, w)
	volumeDisplay := &widget.Label{Text: strconv.Itoa(int(volume)) + "%"}
	volumeSlider := &widget.Slider{Min: 0, Max: 100, Step: 1, Value: float64(volume), Orientation: widget.Horizontal,
		OnChanged: func(f float64) {
			volumeDisplay.SetText(strconv.Itoa(int(f)) + "%")
		},
		OnChangeEnded: func(f float64) {
			showErrorOnFailure(command.SetVolume(uint8(f)), w)
		}}

	volumeMute := &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: func() {
		showErrorOnFailure(command.ToggleVolumeMute(), w)

		muted, err := command.GetVolumeMute()
		showErrorOnFailure(err, w)
		if muted {
			volumeSlider.SetValue(0)
			return
		}

		volume, err := command.GetVolume()
		showErrorOnFailure(err, w)
		volumeSlider.SetValue(float64(volume))
	}}
	volumeUp := &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: func() {
		showErrorOnFailure(command.VolumeUp(), w)

		volume, err := command.GetVolume()
		showErrorOnFailure(err, w)
		volumeSlider.SetValue(float64(volume))
	}}
	volumeDown := &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: func() {
		showErrorOnFailure(command.VolumeDown(), w)

		volume, err := command.GetVolume()
		showErrorOnFailure(err, w)
		volumeSlider.SetValue(float64(volume))
	}}

	inputLabel := &widget.Label{Text: "Select input:", TextStyle: fyne.TextStyle{Bold: true}}

	deviceType := device.H95
	inputs, err := device.GetInputNames(deviceType)
	showErrorOnFailure(err, w)
	source, err := command.GetSourceName(deviceType)
	showErrorOnFailure(err, w)
	inputSelector := &widget.Select{Options: inputs, OnChanged: func(input string) { showErrorOnFailure(command.SetSourceName(deviceType, input), w) }, Selected: source}

	return container.NewVBox(
		power,
		widget.NewSeparator(),
		container.NewVBox(
			container.NewBorder(nil, nil, nil, volumeDisplay, volumeSlider),
			container.NewGridWithColumns(3, volumeMute, volumeDown, volumeUp),
		),
		widget.NewSeparator(),
		inputLabel,
		inputSelector,
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
