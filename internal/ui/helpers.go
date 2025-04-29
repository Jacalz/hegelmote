package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func setEnabled(wid fyne.Disableable, enable bool) {
	if !enable {
		wid.Disable()
		return
	}

	wid.Enable()
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
