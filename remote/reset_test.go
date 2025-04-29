package remote

import (
	"testing"
)

func TestSetResetDelay(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.SetResetDelay(0)
	if err != nil || mock.writeBuf.String() != "-r.0\r" {
		t.Fail()
	}

	mock.FlushToReader()

	_, err = control.SetResetDelay(255)
	if err != nil || mock.writeBuf.String() != "-r.255\r" {
		t.Fail()
	}
}

func TestStopResetDelay(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.StopResetDelay()
	if err != nil || mock.writeBuf.String() != "-r.~\r" {
		t.Fail()
	}
}

func TestGetResetDelay(t *testing.T) {
	control, mock := newControlMock()

	// Set state to stopped.
	control.StopResetDelay()
	mock.FlushToReader()

	delay, err := control.GetResetDelay()
	if err != nil || delay.Minutes != 0 || !delay.Stopped || mock.writeBuf.String() != "-r.?\r" {
		t.Fail()
	}
	mock.Close()

	// Set state delay to 255.
	control.SetResetDelay(255)
	mock.FlushToReader()

	delay, err = control.GetResetDelay()
	if err != nil || delay.Minutes != 255 || delay.Stopped || mock.writeBuf.String() != "-r.?\r" {
		t.Fail()
	}
	mock.Close()
}
