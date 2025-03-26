package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

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
