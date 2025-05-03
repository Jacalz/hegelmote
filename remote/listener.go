package remote

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Jacalz/hegelmote/device"
)

const resetInterval = 2 * time.Minute

// NewControlWithListener returns a controller that listems for state changes from
// the amplifier. This contructor should be used over bare struct setup.
func NewControlWithListener(
	onPower func(bool), onVolume func(Volume),
	onMute func(bool), onInput func(device.Input),
	OnReset func(), OnError func(error),
) *ControlWithListener {
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
	control Control

	OnPowerChange  func(poweredOn bool)
	OnVolumeChange func(volume Volume)
	OnMuteChange   func(muted bool)
	OnInputChange  func(input device.Input)
	OnReset        func()
	OnError        func(err error)

	resetTicker *time.Ticker
	connected   atomic.Bool
	lock        sync.Mutex
}

// GetDeviceType returns the device type of the currently connected amplifier.
func (s *ControlWithListener) GetDeviceType() device.Type {
	return s.control.deviceType
}

// Connect connects to the amplifier and starts the listener.
func (s *ControlWithListener) Connect(host string, model device.Type) error {
	s.connected.Store(true)
	s.resetTicker.Reset(resetInterval)

	defer s.runChangeListener()
	defer s.runResetLoop()

	return s.control.Connect(host, model)
}

// Disconnect disconnects from the amplifier and stops the listener.
func (s *ControlWithListener) Disconnect() error {
	s.connected.Store(false)
	s.resetTicker.Stop()

	s.sendLock()
	defer s.lock.Unlock()

	return s.control.Disconnect()
}

// SetPower sets the amplifier to be on or off depending on the passed bool value.
func (s *ControlWithListener) SetPower(on bool) (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.SetPower(on)
}

// TogglePower toggles between on or off given the current state.
func (s *ControlWithListener) TogglePower() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.TogglePower()
}

// GetPower returns the current power status.
func (s *ControlWithListener) GetPower() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetPower()
}

// SetVolumeMute sets the amplifier to be muted or unmuted given the passed bool value.
func (s *ControlWithListener) SetVolumeMute(muted bool) (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.SetVolumeMute(muted)
}

// ToggleVolumeMute toggles the volume between muted and unmuted given current state.
func (s *ControlWithListener) ToggleVolumeMute() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetVolumeMute()
}

// GetVolumeMute returns the curren state of volume being muted or not.
func (s *ControlWithListener) GetVolumeMute() (bool, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetVolumeMute()
}

// SetVolume sets the volume to the given value.
func (s *ControlWithListener) SetVolume(volume Volume) (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.SetVolume(volume)
}

// VolumeDown decreases the volume one step.
func (s *ControlWithListener) VolumeDown() (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.VolumeDown()
}

// VolumeUp increases the volume one step.
func (s *ControlWithListener) VolumeUp() (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.VolumeUp()
}

// GetVolume returns the current volume value.
func (s *ControlWithListener) GetVolume() (Volume, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetVolume()
}

// SetInput sets the input to the given value.
func (s *ControlWithListener) SetInput(input device.Input) (device.Input, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.SetInput(input)
}

// GetInput returns the currently selected input.
func (s *ControlWithListener) GetInput() (device.Input, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetInput()
}

// SetResetDelay sets a timeout in minutes for when to reset the connection.
func (s *ControlWithListener) SetResetDelay(delay Minutes) (Delay, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.SetResetDelay(delay)
}

// StopResetDelay stops the reset delay from ticking down.
func (s *ControlWithListener) StopResetDelay() (Delay, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.StopResetDelay()
}

// GetResetDelay returns the current delay for reset.
func (s *ControlWithListener) GetResetDelay() (Delay, error) {
	s.sendLock()
	defer s.lock.Unlock()

	return s.control.GetResetDelay()
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (s *ControlWithListener) sendLock() {
	if conn := s.control.conn; conn != nil && conn.SetReadDeadline(time.Now()) == nil {
		defer conn.SetReadDeadline(time.Time{})
	}

	s.lock.Lock()
}

func (s *ControlWithListener) trackState() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	resp, err := s.control.read()
	if err != nil {
		nerr, ok := err.(net.Error)
		if ok && nerr.Timeout() || !s.connected.Load() {
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

			if !s.connected.Load() {
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
