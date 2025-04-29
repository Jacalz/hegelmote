package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-i.0\r")

	_, err := control.SetInput(0)
	assert.Error(t, err)

	mock.Fill("-i.1\r")

	_, err = control.SetInput(1)
	assert.NoError(t, err)
	assert.Equal(t, "-i.1\r", mock.writeBuf.String())

	// Fill reader but clear writer.
	mock.FlushToReader()

	_, err = control.SetInput(8)
	assert.NoError(t, err)
	assert.Equal(t, "-i.8\r", mock.writeBuf.String())
}

func TestGetInput(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-i.1\r")

	_, err := control.SetInput(1)
	assert.NoError(t, err)
	mock.FlushToReader()

	number, err := control.GetInput()
	assert.NoError(t, err)
	assert.Equal(t, 1, number)
	assert.Equal(t, "-i.?\r", mock.writeBuf.String())

	mock.Fill("-i.8\r")
	_, err = control.SetInput(8)
	assert.NoError(t, err)
	mock.FlushToReader()

	number, err = control.GetInput()
	assert.NoError(t, err)
	assert.Equal(t, 8, number)
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
	_, err := control.SetInput(1)
	assert.NoError(t, err)
	mock.FlushToReader()

	name, err := control.GetInputName()
	assert.NoError(t, err)
	assert.Equal(t, "Analog 1", name)

	mock.Close()

	mock.Fill("-i.8\r")
	_, err = control.SetInput(8)
	assert.NoError(t, err)
	mock.FlushToReader()

	name, err = control.GetInputName()
	assert.NoError(t, err)
	assert.Equal(t, "Network", name)

	_, err = control.GetInputName()
	assert.Error(t, err)
}
