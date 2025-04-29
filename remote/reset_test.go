package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetResetDelay(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-r.0\r")

	_, err := control.SetResetDelay(0)
	assert.NoError(t, err)
	assert.Equal(t, "-r.0\r", mock.writeBuf.String())

	mock.Fill("-r.255\r")

	_, err = control.SetResetDelay(255)
	assert.NoError(t, err)
	assert.Equal(t, "-r.255\r", mock.writeBuf.String())
}

func TestStopResetDelay(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-r.0\r")

	_, err := control.StopResetDelay()
	assert.NoError(t, err)
	assert.Equal(t, "-r.~\r", mock.writeBuf.String())
}

func TestGetResetDelay(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-r.~\r")

	_, err := control.StopResetDelay()
	assert.NoError(t, err)
	mock.FlushToReader()

	delay, err := control.GetResetDelay()
	assert.NoError(t, err)
	assert.Zero(t, delay.Minutes)
	assert.True(t, delay.Stopped)
	assert.Equal(t, "-r.?\r", mock.writeBuf.String())

	mock.Fill("-r.255\r")

	_, err = control.SetResetDelay(255)
	assert.NoError(t, err)
	mock.FlushToReader()

	delay, err = control.GetResetDelay()
	assert.NoError(t, err)
	assert.Equal(t, 255, delay.Minutes)
	assert.False(t, delay.Stopped)
	assert.Equal(t, "-r.?\r", mock.writeBuf.String())
}
