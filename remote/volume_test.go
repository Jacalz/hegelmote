package remote

import (
	"testing"
)

func TestSetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.SetVolume(101)
	if err == nil {
		t.Fail()
	}

	err = control.SetVolume(0)
	if err != nil || mock.writeBuf.String() != "-v.0\r" {
		t.Fail()
	}

	mock.FlushToReader()

	err = control.SetVolume(100)
	if err != nil || mock.writeBuf.String() != "-v.100\r" {
		t.Fail()
	}
}

func TestVolumeUp(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.VolumeUp()
	if err != nil || mock.writeBuf.String() != "-v.u\r" {
		t.Fail()
	}
}

func TestVolumeDown(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.VolumeDown()
	if err != nil || mock.writeBuf.String() != "-v.d\r" {
		t.Fail()
	}
}

func TestGetVolume(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

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

func TestSetVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.SetVolumeMute(false)
	if err != nil || mock.writeBuf.String() != "-m.0\r" {
		t.Fail()
	}

	mock.FlushToReader()

	err = control.SetVolumeMute(true)
	if err != nil || mock.writeBuf.String() != "-m.1\r" {
		t.Fail()
	}
}

func TestToggleVolumeMute(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.ToggleVolumeMute()
	if err != nil || mock.writeBuf.String() != "-m.t\r" {
		t.Fail()
	}
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
