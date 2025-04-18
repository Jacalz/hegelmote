package main

import (
	_ "embed"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Jacalz/hegelmote/device"
)

//go:embed assets/img/power.svg
var powerIconContents []byte

type remoteUI struct {
	amplifier statefulController
	current   state
	window    fyne.Window

	// Widgets:
	powerToggle                      *widget.Button
	volumeDisplay                    *widget.Label
	volumeSlider                     *widget.Slider
	volumeMute, volumeDown, volumeUp *widget.Button
	inputLabel                       *widget.Label
	inputSelector                    *widget.Select
}

func (r *remoteUI) refreshPower() {
	if r.current.poweredOn {
		r.powerToggle.SetText("Power off")
	} else {
		r.powerToggle.SetText("Power on")
	}
}

func (r *remoteUI) refreshVolumeSlider() {
	r.volumeSlider.OnChangeEnded = nil

	r.volumeSlider.Value = float64(r.current.volume)
	r.volumeSlider.OnChanged(r.volumeSlider.Value)

	if r.current.poweredOn && !r.current.muted {
		enableAndRefresh(r.volumeSlider)
		setLabelImportance(r.volumeDisplay, widget.MediumImportance)
	} else {
		disableAndRefresh(r.volumeSlider)
		setLabelImportance(r.volumeDisplay, widget.LowImportance)
	}

	r.volumeSlider.OnChangeEnded = r.onVolumeDragEnd
}

func (r *remoteUI) refreshVolumeButtons() {
	if r.current.poweredOn {
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
	r.inputSelector.Selected = r.current.input

	if r.current.poweredOn {
		setLabelImportance(r.inputLabel, widget.MediumImportance)
		enableAndRefresh(r.inputSelector)
	} else {
		setLabelImportance(r.inputLabel, widget.LowImportance)
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
	r.current = r.amplifier.togglePower()
	r.fullRefresh()
}

func (r *remoteUI) onVolumeDrag(percentage float64) {
	r.volumeDisplay.SetText(strconv.Itoa(int(percentage)) + "%")
}

func (r *remoteUI) onVolumeDragEnd(percentage float64) {
	r.current = r.amplifier.setVolume(uint8(percentage))
	r.refreshVolumeSlider()
}

func (r *remoteUI) onMute() {
	r.current = r.amplifier.toggleMute()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onVolumeDown() {
	r.current = r.amplifier.volumeDown()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onVolumeUp() {
	r.current = r.amplifier.volumeUp()
	r.refreshVolumeSlider()
}

func (r *remoteUI) onInputSelect(input string) {
	r.current = r.amplifier.setInput(input)
	r.refreshInput()
}

func (r *remoteUI) connect(host string, model device.Device) error {
	err := r.amplifier.Connect(host, model)
	if err != nil {
		fyne.LogError("Failed to connect to amplifier", err)
		return err
	}

	inputs, err := device.GetInputNames(model)
	if err != nil {
		fyne.LogError("Failed to get input names for model", err)
		return err
	}

	r.inputSelector.Options = inputs
	r.current = r.amplifier.load()
	r.fullRefresh()

	r.amplifier.trackChanges(
		func(refresh refreshed, newState state) {
			fyne.Do(func() {
				r.current = newState

				switch refresh {
				case refreshPower:
					r.fullRefresh()
				case refreshVolume, refreshMute:
					r.refreshVolumeSlider()
				case refreshInput:
					r.refreshInput()
				}
			})
		},
	)

	return nil
}

func (r *remoteUI) setUpConnection(prefs fyne.Preferences, w fyne.Window) {
	host := prefs.String("host")
	modelID := prefs.IntWithFallback("model", -1)
	if host != "" && modelID >= 0 && modelID <= int(device.H590) {
		err := r.connect(host, device.Device(modelID)) // #nosec - Range is checked above!
		if err == nil {
			return
		}

		fyne.LogError("Failed to connect to saved connection", err)
		prefs.RemoveValue("host")
		prefs.RemoveValue("model")
	}

	showConnectionDialog(r, w)
}

func buildRemoteUI(a fyne.App, w fyne.Window) (*remoteUI, fyne.CanvasObject) {
	ui := &remoteUI{window: w}

	powerIcon := theme.NewThemedResource(&fyne.StaticResource{StaticName: "power.svg", StaticContent: powerIconContents})
	ui.powerToggle = &widget.Button{Icon: powerIcon, Text: "Toggle power", OnTapped: ui.onPowerToggle}

	ui.volumeDisplay = &widget.Label{Text: "0%"}
	ui.volumeSlider = &widget.Slider{Min: 0, Max: 100, Step: 1, OnChanged: ui.onVolumeDrag, OnChangeEnded: ui.onVolumeDragEnd}

	ui.volumeMute = &widget.Button{Icon: theme.VolumeMuteIcon(), OnTapped: ui.onMute}

	ui.volumeDown = &widget.Button{Icon: theme.VolumeDownIcon(), OnTapped: ui.onVolumeDown}
	ui.volumeUp = &widget.Button{Icon: theme.VolumeUpIcon(), OnTapped: ui.onVolumeUp}

	ui.inputLabel = &widget.Label{Text: "Select input:", TextStyle: fyne.TextStyle{Bold: true}}
	ui.inputSelector = &widget.Select{PlaceHolder: "Select an input", OnChanged: ui.onInputSelect}

	ui.setUpConnection(a.Preferences(), w)

	return ui, container.NewVBox(
		ui.powerToggle,
		widget.NewSeparator(),
		container.NewVBox(
			container.NewBorder(nil, nil, nil, ui.volumeDisplay, ui.volumeSlider),
			container.NewGridWithColumns(3, ui.volumeMute, ui.volumeDown, ui.volumeUp),
		),
		widget.NewSeparator(),
		ui.inputLabel,
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

func setLabelImportance(label *widget.Label, importance widget.Importance) {
	if label.Importance == importance {
		return
	}

	label.Importance = importance
	label.Refresh()
}
