package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

func buildUI(command *remote.Control) fyne.CanvasObject {
	power := &widget.Button{Text: "Toggle Power", OnTapped: func() { command.TogglePower() }}
	volumeMute := &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: func() { command.ToggleVolumeMute() }}
	volumeUp := &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: func() { command.VolumeUp() }}
	volumeDown := &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: func() { command.VolumeDown() }}

	deviceType := device.H95
	inputs, _ := device.GetInputNames(deviceType)
	inputSelector := &widget.Select{Options: inputs, PlaceHolder: "Select input to use", OnChanged: func(input string) { command.SetSourceName(deviceType, input) }}

	source, _ := command.GetSourceName(device.H95)
	inputSelector.Selected = source

	return container.NewVBox(
		power,
		container.NewGridWithColumns(3, volumeMute, volumeDown, volumeUp),
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

	w.SetContent(buildUI(command))
	w.Resize(fyne.NewSize(300, 400))
	w.SetMaster()
	w.ShowAndRun()
}
