package remote

import (
	"testing"
)

func TestSetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	_, err := control.SetVolume(101)
	if err == nil {
		t.Fail()
	}

	_, err = control.SetVolume(0)
	if err != nil || mock.writeBuf.String() != "-v.0\r" {
		t.Fail()
	}

	mock.Fill("-v.100\r")

	_, err = control.SetVolume(100)
	if err != nil || mock.writeBuf.String() != "-v.100\r" {
		t.Fail()
	}
}

func TestVolumeUp(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.1\r")

	_, err := control.VolumeUp()
	if err != nil || mock.writeBuf.String() != "-v.u\r" {
		t.Fail()
	}
}

func TestVolumeDown(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	_, err := control.VolumeDown()
	if err != nil || mock.writeBuf.String() != "-v.d\r" {
		t.Fail()
	}
}

func TestGetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-v.0\r")

	control.SetVolume(0)
	mock.FlushToReader()

	volume, err := control.GetVolume()
	if err != nil || volume != 0 {
		t.Fail()
	}

	mock.FlushToReader()
	control.SetVolume(100)
	mock.FlushToReader()

	volume, err = control.GetVolume()
	if err != nil || volume != 100 {
		t.Fail()
	}
}
