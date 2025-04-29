package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-m.0\r")

	_, err := control.SetVolumeMute(false)
	assert.NoError(t, err)
	assert.Equal(t, "-m.0\r", mock.writeBuf.String())

	mock.FlushToReader()

	_, err = control.SetVolumeMute(true)
	assert.NoError(t, err)
	assert.Equal(t, "-m.1\r", mock.writeBuf.String())
}

func TestToggleVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-m.1\r")

	muted, err := control.ToggleVolumeMute()
	assert.NoError(t, err)
	assert.True(t, muted)
	assert.Equal(t, "-m.t\r", mock.writeBuf.String())

	mock.Fill("-m.0\r")

	muted, err = control.ToggleVolumeMute()
	assert.NoError(t, err)
	assert.False(t, muted)
	assert.Equal(t, "-m.t\r", mock.writeBuf.String())
}

func TestGetVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-m.0\r")

	_, err := control.SetVolumeMute(false)
	assert.NoError(t, err)
	mock.FlushToReader()

	muted, err := control.GetVolumeMute()
	assert.NoError(t, err)
	assert.False(t, muted)

	mock.Fill("-m.1\r")

	_, err = control.SetVolumeMute(true)
	assert.NoError(t, err)
	mock.FlushToReader()

	muted, err = control.GetVolumeMute()
	assert.NoError(t, err)
	assert.True(t, muted)
}
