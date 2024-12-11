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

	mock.readBuf.WriteString("-p.0\r")
	on, err := control.GetPower()
	if err != nil || on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
	mock.Close()

	mock.readBuf.WriteString("-p.1\r")
	on, err = control.GetPower()
	if err != nil || !on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
}
