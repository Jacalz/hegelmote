package main

import (
	"net"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"github.com/Jacalz/hegelmote/device"
	"github.com/Jacalz/hegelmote/remote"
)

type state struct {
	poweredOn bool
	volume    uint
	muted     bool
	input     string

	control *remote.Control
	lock    sync.Mutex
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (s *state) sendLock() {
	s.control.Conn.SetReadDeadline(time.Now())
	s.lock.Lock()
	s.control.Conn.SetReadDeadline(time.Time{})
}

func (s *state) togglePower() {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.TogglePower()
	if err != nil {
		fyne.LogError("Failed to toggle power", err)
		return
	}

	s.poweredOn = !s.poweredOn
}

func (s *state) setVolume(percentage uint8) {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.SetVolume(percentage)
	if err != nil {
		fyne.LogError("Failed to set volume", err)
		return
	}

	s.volume = uint(percentage)
}

func (s *state) toggleMute() {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.ToggleVolumeMute()
	if err != nil {
		fyne.LogError("Failed to toggle mute", err)
		return
	}

	s.muted = !s.muted
}

func (s *state) volumeDown() {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.VolumeDown()
	if err != nil {
		fyne.LogError("Failed to lower volume", err)
		return
	}

	s.volume = max(0, s.volume-1)
}

func (s *state) volumeUp() {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.VolumeUp()
	if err != nil {
		fyne.LogError("Failed to increase volume", err)
		return
	}

	s.volume = min(100, s.volume+1)
}

func (s *state) setInput(input string) {
	s.sendLock()
	defer s.lock.Unlock()

	err := s.control.SetSourceName(input)
	if err != nil {
		fyne.LogError("Failed to set input", err)
		return
	}

	s.input = input
}

func (s *state) listenForChanges(callback func()) {
	go func() {
		for {
			s.lock.Lock()
			resp, err := s.control.Read()
			s.lock.Unlock()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					continue
				}

				fyne.LogError("Error when listening to changes", err)
				break
			}

			switch resp[1] {
			case 'p':
				s.poweredOn = resp[3] == '1'
			case 'v':
				volume, _ := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
				s.volume = uint(volume)
			case 'm':
				s.muted = resp[3] == '1'
			case 'i':
				input, _ := strconv.ParseUint(string(resp[3:len(resp)-1]), 10, 8)
				s.input, _ = device.NameFromNumber(device.H95, uint(input))
			default:
				continue
			}

			callback()
		}
	}()
}

func (s *state) load() {
	on, err := s.control.GetPower()
	if err != nil {
		fyne.LogError("Failed to read power status", err)
		return
	}

	s.poweredOn = on

	volume, err := s.control.GetVolume()
	if err != nil {
		fyne.LogError("Failed to read volume", err)
		return
	}

	s.volume = volume

	muted, err := s.control.GetVolumeMute()
	if err != nil {
		fyne.LogError("Failed to read mute status", err)
		return
	}

	s.muted = muted

	input, err := s.control.GetSourceName()
	if err != nil {
		fyne.LogError("Failed to get current input", err)
		return
	}

	s.input = input
}
