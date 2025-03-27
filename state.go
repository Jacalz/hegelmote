package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type refreshed uint8

const (
	none refreshed = iota
	refreshPower
	refreshVolume
	refreshMute
	refreshInput
)

type state struct {
	poweredOn bool
	volume    uint
	muted     bool
	input     string
}

type statefulController struct {
	status state

	closing bool
	control *remote.Control
	lock    sync.Mutex
}

func (s *statefulController) disconnect() {
	s.closing = true

	err := s.control.Disconnect()
	if err != nil {
		fyne.LogError("Failure on disconnecting", err)
	}
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (s *statefulController) sendLock() {
	err := s.control.Conn.SetReadDeadline(time.Now())
	if err != nil {
		fyne.LogError("Failure when unblocking state tracker", err)
	}

	s.lock.Lock()

	err = s.control.Conn.SetReadDeadline(time.Time{})
	if err != nil {
		fyne.LogError("Failure when restoring state tracker setup", err)
	}
}

func (s *statefulController) togglePower() state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.TogglePower()
	if err != nil {
		fyne.LogError("Failed to toggle power", err)
		return s.status
	}

	s.status.poweredOn = !s.status.poweredOn
	return s.status
}

func (s *statefulController) setVolume(percentage uint8) state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.SetVolume(percentage)
	if err != nil {
		fyne.LogError("Failed to set volume", err)
		return s.status
	}

	s.status.volume = uint(percentage)
	return s.status
}

func (s *statefulController) toggleMute() state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.ToggleVolumeMute()
	if err != nil {
		fyne.LogError("Failed to toggle mute", err)
		return s.status
	}

	s.status.muted = !s.status.muted
	return s.status
}

func (s *statefulController) volumeDown() state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.VolumeDown()
	if err != nil {
		fyne.LogError("Failed to lower volume", err)
		return s.status
	}

	s.status.volume = max(0, s.status.volume-1)
	return s.status
}

func (s *statefulController) volumeUp() state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.VolumeUp()
	if err != nil {
		fyne.LogError("Failed to increase volume", err)
		return s.status
	}

	s.status.volume = min(100, s.status.volume+1)
	return s.status
}

func (s *statefulController) setInput(input string) state {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.SetSourceName(input)
	if err != nil {
		fyne.LogError("Failed to set input", err)
		return s.status
	}

	s.status.input = input
	return s.status
}

func (s *statefulController) trackState() (refreshed, state, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp, err := s.control.Read()
	if err != nil {
		nerr, ok := err.(net.Error)
		if (ok && nerr.Timeout()) || s.closing {
			return none, s.status, nil
		}

		return none, s.status, err
	}

	switch resp[1] {
	case 'p':
		s.status.poweredOn = resp[3] == '1'
	case 'v':
		volume, err := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
		if err != nil {
			return none, s.status, err
		}
		s.status.volume = uint(volume)
		return refreshVolume, s.status, nil
	case 'm':
		s.status.muted = resp[3] == '1'
		return refreshMute, s.status, nil
	case 'i':
		input, err := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
		if err != nil {
			return none, s.status, err
		}

		inputName, err := device.NameFromNumber(device.H95, uint(input))
		if err != nil {
			return none, s.status, err
		}
		s.status.input = inputName
		return refreshInput, s.status, nil
	case 'e':
		return none, s.status, fmt.Errorf("got error code %d from amplifier", resp[3])
	}

	return none, s.status, fmt.Errorf("unknown command \"%c\" received from amplifier", resp[1])
}

func (s *statefulController) trackChanges(callback func(refreshed, state)) {
	go func() {
		for {
			refresh, status, err := s.trackState()
			if err != nil {
				fyne.LogError("Error on tracking state change from amplifier", err)
				return
			}

			if refresh != none {
				callback(refresh, status)
			}
		}
	}()
}

func (s *statefulController) load() state {
	on, err := s.control.GetPower()
	if err != nil {
		fyne.LogError("Failed to read power status", err)
		return s.status
	}

	s.status.poweredOn = on

	volume, err := s.control.GetVolume()
	if err != nil {
		fyne.LogError("Failed to read volume", err)
		return s.status
	}

	s.status.volume = volume

	muted, err := s.control.GetVolumeMute()
	if err != nil {
		fyne.LogError("Failed to read mute status", err)
		return s.status
	}

	s.status.muted = muted

	input, err := s.control.GetSourceName()
	if err != nil {
		fyne.LogError("Failed to get current input", err)
		return s.status
	}

	s.status.input = input
	return s.status
}
