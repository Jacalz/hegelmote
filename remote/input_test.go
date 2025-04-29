package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-i.0\r")

	_, err := control.SetInput(0)
	if err == nil {
		t.Fail()
	}

	mock.Fill("-i.1\r")

	_, err = control.SetInput(1)
	if err != nil || mock.writeBuf.String() != "-i.1\r" {
		t.Fail()
	}

	// Fill reader but clear writer.
	mock.FlushToReader()

	_, err = control.SetInput(8)
	if err != nil || mock.writeBuf.String() != "-i.8\r" {
		t.Fail()
	}
}

func TestGetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-i.1\r")

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

	mock.Fill("-i.1\r")

	_, err := control.SetInputFromName("Analog 1")
	assert.NoError(t, err)
	assert.Equal(t, "-i.1\r", mock.writeBuf.String())

	mock.FlushToReader()

	number, err := control.GetInput()
	assert.NoError(t, err)
	assert.Equal(t, 1, number)

	mock.Fill("-i.8\r")

	_, err = control.SetInputFromName("Network")
	assert.NoError(t, err)
	assert.Equal(t, "-i.8\r", mock.writeBuf.String())

	mock.FlushToReader()

	number, err = control.GetInput()
	assert.NoError(t, err)
	assert.Equal(t, 8, number)

	_, err = control.SetInputFromName("Bogus")
	assert.Error(t, err)
}

func TestGetInputName(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-i.1\r")
	control.SetInput(1)
	mock.FlushToReader()

	name, err := control.GetInputName()
	if err != nil || name != "Analog 1" {
		t.Fail()
	}

	mock.Close()

	mock.Fill("-i.8\r")
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
