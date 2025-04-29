package remote

import (
	"fmt"
	"testing"
)

func TestSetPower(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.SetPower(true)
	if err != nil || mock.writeBuf.String() != "-p.1\r" {
		fmt.Println(err)
		t.Fail()
	}

	mock.FlushToReader()

	_, err = control.SetPower(false)
	if err != nil || mock.writeBuf.String() != "-p.0\r" {
		t.Fail()
	}
}

func TestTogglePower(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.TogglePower()
	if err != nil || mock.writeBuf.String() != "-p.t\r" {
		t.Fail()
	}
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
