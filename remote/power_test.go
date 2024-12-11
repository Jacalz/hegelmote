package remote

import "testing"

func TestSetPower(t *testing.T) {
	control, mock := newControlMock()

	err := control.SetPower(true)
	if err != nil || mock.writeBuf.String() != "-p.1\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()

	err = control.SetPower(false)
	if err != nil || mock.writeBuf.String() != "-p.0\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()
}

func TestTogglePower(t *testing.T) {
	control, mock := newControlMock()

	err := control.TogglePower()
	if err != nil || mock.writeBuf.String() != "-p.t\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()
}

func TestGetPower(t *testing.T) {
	control, mock := newControlMock()

	control.SetPower(false)
	mock.FlushToReader()

	on, err := control.GetPower()
	if err != nil || on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
	mock.Close()

	control.SetPower(true)
	mock.FlushToReader()

	on, err = control.GetPower()
	if err != nil || !on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
}
