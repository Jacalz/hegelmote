package remote

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Jacalz/hegelmote/device"
)

const resetInterval = 2 * time.Minute

// NewControlWithListener returns a controller that listems for state changes from
// the amplifier. This contructor should be used over bare struct setup.
func NewControlWithListener(
	onPower func(bool), onVolume func(Volume),
	onMute func(bool), onInput func(device.Input),
	OnReset func(), OnError func(error)) *ControlWithListener {
	return &ControlWithListener{
		resetTicker:    time.NewTicker(resetInterval),
		OnPowerChange:  onPower,
		OnVolumeChange: onVolume,
		OnMuteChange:   onMute,
		OnInputChange:  onInput,
		OnReset:        OnReset,
		OnError:        OnError,
	}
}

// ControlWithListener is a remote control that listens for changes
// sent from the amplifier. It also sends a reset to the amplifier with
// a fixed delay to allow reconnecting in case of error.
// This data type is thread safe.
type ControlWithListener struct {
	Control

	OnPowerChange  func(poweredOn bool)
	OnVolumeChange func(volume Volume)
	OnMuteChange   func(muted bool)
	OnInputChange  func(input device.Input)
	OnReset        func()
	OnError        func(err error)

	resetTicker *time.Ticker
	closing     bool
	lock        sync.Mutex
}

// Connect tries to connect to the amplifier.
func (s *ControlWithListener) Connect(host string, model device.Device) error {
	s.closing = false
	s.resetTicker.Reset(resetInterval)

	defer s.runChangeListener()
	defer s.runResetLoop()

	return s.Control.Connect(host, model)
}

// Disconnect disconnects from the amplifier.
func (s *ControlWithListener) Disconnect() error {
	s.sendLock()
	defer s.lock.Unlock()

	s.closing = true
	s.resetTicker.Stop()

	return s.Control.Disconnect()
}

// TogglePower sends a toggle command to the amplifier.
func (s *ControlWithListener) TogglePower() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.TogglePower()
}

// SetVolume sets the volume to the given value.
func (s *ControlWithListener) SetVolume(volume Volume) (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.SetVolume(volume)
}

// ToggleMute toggles the mute value.
func (s *ControlWithListener) ToggleMute() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.ToggleVolumeMute()
}

// VolumeDown decreases the volume one step.
func (s *ControlWithListener) VolumeDown() (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.VolumeDown()
}

// VolumeUp increases the volume one step.
func (s *ControlWithListener) VolumeUp() (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.VolumeUp()
}

// SetInput sets the input to the given value.
func (s *ControlWithListener) SetInput(input device.Input) (device.Input, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.SetInput(input)
}

// SetResetDelay sends a reset delay to the amplifier.
func (s *ControlWithListener) SetResetDelay(delay Minutes) (Delay, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.Control.SetResetDelay(delay)
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (s *ControlWithListener) sendLock() {
	if s.conn != nil && s.conn.SetReadDeadline(time.Now()) == nil {
		defer s.conn.SetReadDeadline(time.Time{})
	}

	s.lock.Lock()
}
func (s *ControlWithListener) trackState() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp, err := s.read()
	if err != nil {
		nerr, ok := err.(net.Error)
		if ok && nerr.Timeout() {
			return nil
		}

		if s.closing {
			return nil
		}

		return err
	}

	switch resp[1] {
	case 'p':
		s.OnPowerChange(resp[3] == '1')
	case 'v':
		volume, err := parseUint8FromBuf(resp)
		if err != nil {
			return err
		}
		s.OnVolumeChange(volume)
	case 'm':
		s.OnMuteChange(resp[3] == '1')
	case 'i':
		input, err := parseUint8FromBuf(resp)
		if err != nil {
			return err
		}
		s.OnInputChange(input)
	case 'r':
		if resp[3] == '0' {
			s.OnReset()
		}
	case 'e':
		return errorFromCode(resp[3])
	default:
		return fmt.Errorf("received unknown command \"%c\" from amplifier", resp[1])
	}

	return nil
}

func (s *ControlWithListener) runChangeListener() {
	go func() {
		for {
			err := s.trackState()
			if err != nil {
				s.OnError(err)
				return
			}

			if s.closing {
				return
			}
		}
	}()
}

func (s *ControlWithListener) runResetLoop() {
	_, err := s.SetResetDelay(3)
	if err != nil {
		s.OnError(err)
	}

	go func() {
		for range s.resetTicker.C {
			_, err := s.SetResetDelay(3)
			if err != nil {
				s.OnError(err)
			}
		}
	}()
}
