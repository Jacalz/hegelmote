package remote

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/Jacalz/hegelmote/device"
)

var errConnectionReset = errors.New("connection was reset")

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

	c.conn.reads = make(chan readResponse)
	c.resetTicker.Stop()
	go c.runResetLoop()
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
	sending     atomic.Bool
	closing     atomic.Bool
	conn        listenerConn
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

	// Update connection wrapper with new connection:
	c.conn.Conn = c.control.conn
	c.control.conn = &c.conn

	c.closing.Store(false)
	go c.runChangeListener()

	c.resetTicker.Reset(resetInterval)
	_, err = c.SetResetDelay(3)
	return err
}

// Disconnect disconnects from the amplifier and stops the listener.
func (c *ControlWithListener) Disconnect() error {
	c.resetTicker.Stop()
	c.closing.Store(true)
	err := c.control.Disconnect()
	c.conn.Conn = nil
	return err
}

// SetPower sets the amplifier to be on or off depending on the passed bool value.
func (c *ControlWithListener) SetPower(on bool) (bool, error) {
	return sendWithArgument(&c.sending, c.control.SetPower, on)
}

// TogglePower toggles between on or off given the current state.
func (c *ControlWithListener) TogglePower() (bool, error) {
	return send(&c.sending, c.control.TogglePower)
}

// GetPower returns the current power status.
func (c *ControlWithListener) GetPower() (bool, error) {
	return send(&c.sending, c.control.GetPower)
}

// SetVolumeMute sets the amplifier to be muted or unmuted given the passed bool value.
func (c *ControlWithListener) SetVolumeMute(muted bool) (bool, error) {
	return sendWithArgument(&c.sending, c.control.SetVolumeMute, muted)
}

// ToggleVolumeMute toggles the volume between muted and unmuted given current state.
func (c *ControlWithListener) ToggleVolumeMute() (bool, error) {
	return send(&c.sending, c.control.ToggleVolumeMute)
}

// GetVolumeMute returns the curren state of volume being muted or not.
func (c *ControlWithListener) GetVolumeMute() (bool, error) {
	return send(&c.sending, c.control.GetVolumeMute)
}

// SetVolume sets the volume to the given value.
func (c *ControlWithListener) SetVolume(volume Volume) (Volume, error) {
	return sendWithArgument(&c.sending, c.control.SetVolume, volume)
}

// VolumeDown decreases the volume one step.
func (c *ControlWithListener) VolumeDown() (Volume, error) {
	return send(&c.sending, c.control.VolumeDown)
}

// VolumeUp increases the volume one step.
func (c *ControlWithListener) VolumeUp() (Volume, error) {
	return send(&c.sending, c.control.VolumeUp)
}

// GetVolume returns the current volume value.
func (c *ControlWithListener) GetVolume() (Volume, error) {
	return send(&c.sending, c.control.GetVolume)
}

// SetInput sets the input to the given value.
func (c *ControlWithListener) SetInput(input device.Input) (device.Input, error) {
	return sendWithArgument(&c.sending, c.control.SetInput, input)
}

// GetInput returns the currently selected input.
func (c *ControlWithListener) GetInput() (device.Input, error) {
	return send(&c.sending, c.control.GetInput)
}

// SetResetDelay sets a timeout in minutes for when to reset the connection.
func (c *ControlWithListener) SetResetDelay(delay Minutes) (Delay, error) {
	return sendWithArgument(&c.sending, c.control.SetResetDelay, delay)
}

// StopResetDelay stops the reset delay from ticking down.
func (c *ControlWithListener) StopResetDelay() (Delay, error) {
	return send(&c.sending, c.control.StopResetDelay)
}

// GetResetDelay returns the current delay for reset.
func (c *ControlWithListener) GetResetDelay() (Delay, error) {
	return send(&c.sending, c.control.GetResetDelay)
}

func (c *ControlWithListener) waitForResponse() error {
	buf := [len("-v.100\r")]byte{}
	n, err := c.conn.Conn.Read(buf[:])
	if c.sending.CompareAndSwap(true, false) {
		c.conn.reads <- readResponse{n: n, buf: buf[:], err: err}
		return nil
	} else if err != nil {
		return err
	}

	resp, err := c.control.verifyResponse(buf, n)
	if err != nil {
		return err
	}

	switch resp[1] {
	case 'p':
		c.OnPowerChange(resp[3] == '1')
	case 'm':
		c.OnMuteChange(resp[3] == '1')
	case 'v', 'i':
		number, err := parseUint8FromBuf(resp)
		if err != nil {
			return err
		}

		if resp[1] == 'v' {
			c.OnVolumeChange(number)
		} else {
			c.OnInputChange(number)
		}
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
	for {
		err := c.waitForResponse()
		if err != nil {
			if !errors.Is(err, errConnectionReset) && !c.closing.Load() {
				c.OnError(err)
			}
			return
		}
	}
}

func (c *ControlWithListener) runResetLoop() {
	for range c.resetTicker.C {
		_, err := c.SetResetDelay(3)
		if err != nil {
			c.OnError(err)
		}
	}
}

type readResponse struct {
	buf []byte
	n   int
	err error
}

type listenerConn struct {
	net.Conn
	reads chan readResponse
}

func (l *listenerConn) Read(p []byte) (int, error) {
	got := <-l.reads
	copy(p, got.buf)
	return got.n, got.err
}

func send[T any](sending *atomic.Bool, cmd func() (T, error)) (T, error) {
	sending.Store(true)
	return cmd()
}

func sendWithArgument[T, U any](sending *atomic.Bool, cmd func(U) (T, error), input U) (T, error) {
	sending.Store(true)
	return cmd(input)
}
