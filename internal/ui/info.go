package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/hegelmote/device"
)

func (m *mainUI) onConnectionInfo() {
	info := &widget.Form{Items: []*widget.FormItem{
		{Text: "Address", Widget: &widget.Label{Text: m.host}},
		{Text: "Model", Widget: &widget.Label{Text: "Hegel " + device.SupportedDeviceNames()[m.amplifier.GetModel()]}},
		{Text: "Status", Widget: &widget.Label{Text: m.connectionLabel.Text}},
	}}

	prefs := fyne.CurrentApp().Preferences()
	var infoDialog *dialog.CustomDialog

	disconnect := &widget.Button{Text: "Disconnect", Icon: theme.CancelIcon(), Importance: widget.LowImportance, OnTapped: func() {
		infoDialog.Hide()
		m.Disconnect()
		showConnectionDialog(m, m.window)
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
