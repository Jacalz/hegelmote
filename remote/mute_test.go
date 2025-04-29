package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-m.0\r")

	_, err := control.SetVolumeMute(false)
	if err != nil || mock.writeBuf.String() != "-m.0\r" {
		t.Fail()
	}

	mock.FlushToReader()

	_, err = control.SetVolumeMute(true)
	if err != nil || mock.writeBuf.String() != "-m.1\r" {
		t.Fail()
	}
}

func TestToggleVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-m.1\r")

	muted, err := control.ToggleVolumeMute()
	assert.NoError(t, err)
	assert.True(t, muted)
	assert.Equal(t, "-m.t\r", mock.writeBuf.String())

	mock.Fill("-m.0\r")
	mock.writeBuf.Reset()

	muted, err = control.ToggleVolumeMute()
	assert.NoError(t, err)
	assert.False(t, muted)
	assert.Equal(t, "-m.t\r", mock.writeBuf.String())
}

func TestGetVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	control.SetVolumeMute(false)
	mock.FlushToReader()

	muted, err := control.GetVolumeMute()
	if err != nil || muted {
		t.Fail()
	}

	mock.Close()

	control.SetVolumeMute(true)
	mock.FlushToReader()

	muted, err = control.GetVolumeMute()
	if err != nil || !muted {
		t.Fail()
	}
}
