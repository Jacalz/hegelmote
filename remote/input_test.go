package remote

import (
	"fmt"
	"testing"
)

func TestSetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.SetInput(0)
	if err == nil {
		t.Fail()
	}

	// Command returns the currently set input on success. Fill buffer.
	mock.readBuf.WriteString("-i.1\r")

	_, err = control.SetInput(1)
	if err != nil || mock.writeBuf.String() != "-i.1\r" {
		t.Fail()
	}

	// Fill reader but clear writer.
	mock.FlushToReader()

	_, err = control.SetInput(8)
	if err != nil || mock.writeBuf.String() != "-i.8\r" {
		fmt.Println(err, mock.writeBuf.String())
		t.Fail()
	}
}

func TestGetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	control.SetInput(1)
	mock.FlushToReader()

	number, err := control.GetInput()
	if err != nil || number != 1 || mock.writeBuf.String() != "-i.?\r" {
		t.Fail()
	}

	mock.FlushToReader()
	control.SetInput(8)
	mock.FlushToReader()

	number, err = control.GetInput()
	if err != nil || number != 8 {
		t.Fail()
	}
}

func TestSetInputFromName(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	_, err := control.SetInputFromName("Analog 1")
	if err != nil || mock.writeBuf.String() != "-i.1\r" {
		t.Fail()
	}

	mock.FlushToReader()

	number, err := control.GetInput()
	if err != nil || number != 1 {
		t.Fail()
	}

	mock.FlushToReader()

	_, err = control.SetInputFromName("Network")
	if err != nil || mock.writeBuf.String() != "-i.8\r" {
		t.Fail()
	}

	mock.FlushToReader()

	number, err = control.GetInput()
	if err != nil || number != 8 {
		t.Fail()
	}

	_, err = control.SetInputFromName("Bogus")
	if err == nil {
		t.Fail()
	}
}

func TestGetInputName(t *testing.T) {
	control, mock := newControlMock()

	mock.readBuf.WriteString("-i.1\r")
	control.SetInput(1)
	mock.FlushToReader()

	name, err := control.GetInputName()
	if err != nil || name != "Analog 1" {
		t.Fail()
	}

	mock.Close()

	mock.readBuf.WriteString("-i.8\r")
	control.SetInput(8)
	mock.FlushToReader()

	name, err = control.GetInputName()
	if err != nil || name != "Network" {
		t.Fail()
	}

	_, err = control.GetInputName()
	if err == nil {
		t.Fail()
	}
}
