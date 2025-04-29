package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	_, err := control.SetVolume(101)
	assert.Error(t, err)

	_, err = control.SetVolume(0)
	assert.NoError(t, err)
	assert.Equal(t, "-v.0\r", mock.writeBuf.String())

	mock.Fill("-v.50\r")

	_, err = control.SetVolume(50)
	assert.NoError(t, err)
	assert.Equal(t, "-v.50\r", mock.writeBuf.String())

	mock.Fill("-v.100\r")

	_, err = control.SetVolume(100)
	assert.NoError(t, err)
	assert.Equal(t, "-v.100\r", mock.writeBuf.String())
}

func TestVolumeUp(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.1\r")

	_, err := control.VolumeUp()
	assert.NoError(t, err)
	assert.Equal(t, "-v.u\r", mock.writeBuf.String())
}

func TestVolumeDown(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	_, err := control.VolumeDown()
	assert.NoError(t, err)
	assert.Equal(t, "-v.d\r", mock.writeBuf.String())
}

func TestGetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	_, err := control.SetVolume(0)
	assert.NoError(t, err)
	mock.FlushToReader()

	volume, err := control.GetVolume()
	assert.NoError(t, err)
	assert.Zero(t, volume)

	mock.Fill("-v.100\r")
	_, err = control.SetVolume(100)
	assert.NoError(t, err)
	mock.FlushToReader()

	volume, err = control.GetVolume()
	assert.NoError(t, err)
	assert.Equal(t, 100, volume)
}
