package remote

import "testing"

func TestSetPower(t *testing.T) {
	control, mock := newControlMock()

	control.SetPower(true)
	if mock.writeBuf.String() != "-p.1\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()

	control.SetPower(false)
	if mock.writeBuf.String() != "-p.0\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()
}

func TestTogglePower(t *testing.T) {
	control, mock := newControlMock()

	control.TogglePower()
	if mock.writeBuf.String() != "-p.t\r" {
		t.Fail()
	}
	mock.writeBuf.Reset()
}

func TestGetPower(t *testing.T) {
	control, mock := newControlMock()

	mock.readBuf.WriteString("-p.0\r")
	on, _ := control.GetPower()
	if on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
	mock.Close()

	mock.readBuf.WriteString("-p.1\r")
	on, _ = control.GetPower()
	if !on || mock.writeBuf.String() != "-p.?\r" {
		t.Fail()
	}
}
