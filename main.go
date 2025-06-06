package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Jacalz/hegelmote/internal/ui"
)

func main() {
	stop := profile()
	if stop != nil {
		defer stop()
	}

	a := app.NewWithID("io.github.jacalz.hegelmote")
	w := a.NewWindow("Hegelmote")

	ui, content := ui.Build(a, w)
	defer ui.Disconnect()

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 400))
	w.SetMaster()
	w.ShowAndRun()
}
