package ui

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

const resetInterval = 2 * time.Minute

type refreshed uint8

const (
	none refreshed = iota
	closing
	refreshPower
	refreshVolume
	refreshMute
	refreshInput
	reset
)

type state struct {
	poweredOn bool
	volume    remote.Volume
	muted     bool
	input     device.Input
}

type statefulController struct {
	remote.Control

	resetTicker *time.Ticker
	closing     bool
	lock        sync.Mutex
}

func (s *statefulController) disconnect() {
	s.sendLock()
	defer s.lock.Unlock()

	s.closing = true
	if s.resetTicker != nil {
		s.resetTicker.Stop()
	}

	err := s.Disconnect()
	if err != nil {
		fyne.LogError("Failure on disconnecting", err)
	}
}

func (s *statefulController) runResetLoop() {
	if s.resetTicker == nil {
		s.resetTicker = time.NewTicker(resetInterval)
	} else {
		s.resetTicker.Reset(resetInterval)
	}

	go func() {
		s.reset(3)
		for range s.resetTicker.C {
			s.reset(3)
		}
	}()
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (s *statefulController) sendLock() {
	if s.Conn != nil && s.Conn.SetReadDeadline(time.Now()) == nil {
		defer s.Conn.SetReadDeadline(time.Time{})
	}

	s.lock.Lock()
}

func (s *statefulController) togglePower() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.TogglePower()
}

func (s *statefulController) setVolume(volume remote.Volume) (remote.Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.SetVolume(volume)
}

func (s *statefulController) toggleMute() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.ToggleVolumeMute()
}

func (s *statefulController) volumeDown() (remote.Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.VolumeDown()
}

func (s *statefulController) volumeUp() (remote.Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.VolumeUp()
}

func (s *statefulController) setInput(input device.Input) (device.Input, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.SetInput(input)
}

func (s *statefulController) reset(delay remote.Minutes) (remote.Delay, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.SetResetDelay(delay)
}

func (s *statefulController) trackState() (refreshed, state, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	status := state{}

	resp, err := s.Read()
	if err != nil {
		nerr, ok := err.(net.Error)
		if ok && nerr.Timeout() {
			return none, status, nil
		}

		if s.closing {
			return closing, status, nil
		}

		return none, status, err
	}

	switch resp[1] {
	case 'p':
		status.poweredOn = resp[3] == '1'
		return refreshPower, status, nil
	case 'v':
		volume, err := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
		if err != nil {
			return none, status, err
		}
		status.volume = remote.Volume(volume)
		return refreshVolume, status, nil
	case 'm':
		status.muted = resp[3] == '1'
		return refreshMute, status, nil
	case 'i':
		input, err := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
		if err != nil {
			return none, status, err
		}
		status.input = device.Input(input)
		return refreshInput, status, nil
	case 'r':
		if resp[3] == '0' {
			return reset, status, nil
		}
		return none, status, nil
	case 'e':
		return none, status, fmt.Errorf("got error code %d from amplifier", resp[3])
	}

	return none, status, fmt.Errorf("unknown command \"%c\" received from amplifier", resp[1])
}

func (s *statefulController) trackChanges(callback func(refreshed, state)) {
	go func() {
		for {
			refresh, status, err := s.trackState()
			if err != nil {
				fyne.LogError("Error on tracking state change from amplifier", err)
				return
			}

			if refresh == closing {
				return
			} else if refresh != none {
				callback(refresh, status)
			}
		}
	}()
}
