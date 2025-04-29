package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetPower(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-p.1\r")

	_, err := control.SetPower(true)
	assert.NoError(t, err)
	assert.Equal(t, "-p.1\r", mock.writeBuf.String())

	mock.Fill("-p.0\r")

	_, err = control.SetPower(false)
	assert.NoError(t, err)
	assert.Equal(t, "-p.0\r", mock.writeBuf.String())
}

func TestTogglePower(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-p.0\r")

	_, err := control.TogglePower()
	assert.NoError(t, err)
	assert.Equal(t, "-p.t\r", mock.writeBuf.String())
}

func TestGetPower(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-p.0\r")

	_, err := control.SetPower(false)
	assert.NoError(t, err)
	mock.FlushToReader()

	on, err := control.GetPower()
	assert.NoError(t, err)
	assert.False(t, on)
	assert.Equal(t, "-p.?\r", mock.writeBuf.String())

	mock.Fill("-p.1\r")

	_, err = control.SetPower(true)
	assert.NoError(t, err)
	mock.FlushToReader()

	on, err = control.GetPower()
	assert.NoError(t, err)
	assert.True(t, on)
	assert.Equal(t, "-p.?\r", mock.writeBuf.String())
}
