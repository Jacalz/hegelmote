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
	c := &ControlWithListener{
		resetTicker:    time.NewTicker(resetInterval),
		OnPowerChange:  onPower,
		OnVolumeChange: onVolume,
		OnMuteChange:   onMute,
		OnInputChange:  onInput,
		OnReset:        OnReset,
		OnError:        OnError,
	}

	c.resetTicker.Stop()
	c.runResetLoop()
	return c
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
func (c *ControlWithListener) GetDeviceType() device.Type {
	return c.control.deviceType
}

// Connect connects to the amplifier and starts the listener.
func (c *ControlWithListener) Connect(host string, model device.Type) error {
	err := c.control.Connect(host, model)
	if err != nil {
		return err
	}

	c.connected.Store(true)
	c.runChangeListener()

	c.resetTicker.Reset(resetInterval)
	_, err = c.SetResetDelay(3)
	return err
}

// Disconnect disconnects from the amplifier and stops the listener.
func (c *ControlWithListener) Disconnect() error {
	c.connected.Store(false)
	c.resetTicker.Stop()

	c.sendLock()
	defer c.lock.Unlock()

	return c.control.Disconnect()
}

// SetPower sets the amplifier to be on or off depending on the passed bool value.
func (c *ControlWithListener) SetPower(on bool) (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.SetPower(on)
}

// TogglePower toggles between on or off given the current state.
func (c *ControlWithListener) TogglePower() (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.TogglePower()
}

// GetPower returns the current power statuc.
func (c *ControlWithListener) GetPower() (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.GetPower()
}

// SetVolumeMute sets the amplifier to be muted or unmuted given the passed bool value.
func (c *ControlWithListener) SetVolumeMute(muted bool) (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.SetVolumeMute(muted)
}

// ToggleVolumeMute toggles the volume between muted and unmuted given current state.
func (c *ControlWithListener) ToggleVolumeMute() (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.ToggleVolumeMute()
}

// GetVolumeMute returns the curren state of volume being muted or not.
func (c *ControlWithListener) GetVolumeMute() (bool, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.GetVolumeMute()
}

// SetVolume sets the volume to the given value.
func (c *ControlWithListener) SetVolume(volume Volume) (Volume, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.SetVolume(volume)
}

// VolumeDown decreases the volume one step.
func (c *ControlWithListener) VolumeDown() (Volume, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.VolumeDown()
}

// VolumeUp increases the volume one step.
func (c *ControlWithListener) VolumeUp() (Volume, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.VolumeUp()
}

// GetVolume returns the current volume value.
func (c *ControlWithListener) GetVolume() (Volume, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.GetVolume()
}

// SetInput sets the input to the given value.
func (c *ControlWithListener) SetInput(input device.Input) (device.Input, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.SetInput(input)
}

// GetInput returns the currently selected input.
func (c *ControlWithListener) GetInput() (device.Input, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.GetInput()
}

// SetResetDelay sets a timeout in minutes for when to reset the connection.
func (c *ControlWithListener) SetResetDelay(delay Minutes) (Delay, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.SetResetDelay(delay)
}

// StopResetDelay stops the reset delay from ticking down.
func (c *ControlWithListener) StopResetDelay() (Delay, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.StopResetDelay()
}

// GetResetDelay returns the current delay for reset.
func (c *ControlWithListener) GetResetDelay() (Delay, error) {
	c.sendLock()
	defer c.lock.Unlock()

	return c.control.GetResetDelay()
}

// sendLock unblocks the reading state tracker, locks and reverts back to blocking read.
func (c *ControlWithListener) sendLock() {
	if conn := c.control.conn; conn != nil && conn.SetReadDeadline(time.Now()) == nil {
		defer conn.SetReadDeadline(time.Time{})
	}

	c.lock.Lock()
}

func (c *ControlWithListener) trackState() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.control.read()
	if err != nil {
		nerr, ok := err.(net.Error)
		if ok && nerr.Timeout() || !c.connected.Load() {
			return nil
		}

		return err
	}

	switch resp[1] {
	case 'p':
		c.OnPowerChange(resp[3] == '1')
	case 'v':
		volume, err := parseUint8FromBuf(resp)
		if err != nil {
			return err
		}
		c.OnVolumeChange(volume)
	case 'm':
		c.OnMuteChange(resp[3] == '1')
	case 'i':
		input, err := parseUint8FromBuf(resp)
		if err != nil {
			return err
		}
		c.OnInputChange(input)
	case 'r':
		if resp[3] == '0' {
			c.OnReset()
		}
	case 'e':
		return errorFromCode(resp[3])
	default:
		return fmt.Errorf("received unknown command \"%c\" from amplifier", resp[1])
	}

	return nil
}

func (c *ControlWithListener) runChangeListener() {
	go func() {
		for {
			err := c.trackState()
			if err != nil {
				c.OnError(err)
				return
			}

			if !c.connected.Load() {
				return
			}
		}
	}()
}

func (c *ControlWithListener) runResetLoop() {
	go func() {
		for range c.resetTicker.C {
			_, err := c.SetResetDelay(3)
			if err != nil {
				c.OnError(err)
			}
		}
	}()
}
