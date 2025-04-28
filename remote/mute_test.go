package remote

import "testing"

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
