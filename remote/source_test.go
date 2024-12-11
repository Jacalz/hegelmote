package remote

import (
	"fmt"
	"testing"

	"github.com/Jacalz/hegelmote/device"
)

func TestSetSourceNumber(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.SetSourceNumber(0)
	if err == nil {
		t.Fail()
	}

	// Command returns the currently set input on success. Fill buffer.
	mock.readBuf.WriteString("-i.1\r")

	err = control.SetSourceNumber(1)
	if err != nil || mock.writeBuf.String() != "-i.1\r" {
		t.Fail()
	}

	// Fill reader but clear writer.
	mock.FlushToReader()

	err = control.SetSourceNumber(8)
	if err != nil || mock.writeBuf.String() != "-i.8\r" {
		fmt.Println(err, mock.writeBuf.String())
		t.Fail()
	}
}

func TestGetSourceNumber(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	control.SetSourceNumber(1)
	mock.FlushToReader()

	number, err := control.GetSourceNumber()
	if err != nil || number != 1 || mock.writeBuf.String() != "-i.?\r" {
		t.Fail()
	}

	mock.FlushToReader()
	control.SetSourceNumber(8)
	mock.FlushToReader()

	number, err = control.GetSourceNumber()
	if err != nil || number != 8 {
		t.Fail()
	}
}

func TestSetSouceName(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill()

	err := control.SetSourceName(device.H95, "Analog 1")
	if err != nil || mock.writeBuf.String() != "-i.1\r" {
		t.Fail()
	}

	mock.FlushToReader()

	number, err := control.GetSourceNumber()
	if err != nil || number != 1 {
		t.Fail()
	}

	mock.FlushToReader()

	err = control.SetSourceName(device.H95, "Network")
	if err != nil || mock.writeBuf.String() != "-i.8\r" {
		t.Fail()
	}

	mock.FlushToReader()

	number, err = control.GetSourceNumber()
	if err != nil || number != 8 {
		t.Fail()
	}

	err = control.SetSourceName(device.H590, "Bogus")
	if err == nil {
		t.Fail()
	}
}

func TestGetSourceName(t *testing.T) {
	control, mock := newControlMock()

	mock.readBuf.WriteString("-i.1\r")
	control.SetSourceNumber(1)
	mock.FlushToReader()

	name, err := control.GetSourceName(device.H95)
	if err != nil || name != "Analog 1" {
		t.Fail()
	}

	mock.Close()

	mock.readBuf.WriteString("-i.8\r")
	control.SetSourceNumber(8)
	mock.FlushToReader()

	name, err = control.GetSourceName(device.H95)
	if err != nil || name != "Network" {
		t.Fail()
	}

	_, err = control.GetSourceName(15)
	if err == nil {
		t.Fail()
	}
}
