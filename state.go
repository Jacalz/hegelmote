package main

import (
	"fyne.io/fyne/v2"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type state struct {
	poweredOn bool
	volume    uint
	muted     bool
	inputs    []string
	input     string
}

func readState(control *remote.Control, model device.Device) *state {
	on, err := control.GetPower()
	if err != nil {
		fyne.LogError("Failed to read power status", err)
		return nil
	}

	// Volume:
	volume, err := control.GetVolume()
	if err != nil {
		fyne.LogError("Failed to read volume", err)
		return nil
	}

	// Mute:
	muted, err := control.GetVolumeMute()
	if err != nil {
		fyne.LogError("Failed to read mute status", err)
		return nil
	}

	// Input:
	inputs, err := device.GetInputNames(model)
	if err != nil {
		fyne.LogError("Failed to get input names for device", err)
		return nil
	}

	input, err := control.GetSourceName(model)
	if err != nil {
		fyne.LogError("Failed to get current input", err)
		return nil
	}

	return &state{
		poweredOn: on,
		volume:    volume,
		muted:     muted,
		inputs:    inputs,
		input:     input,
	}
}
