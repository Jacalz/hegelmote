package remote

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestReturnedErrors(t *testing.T) {
	control, mock := newControlMock()

	mock.Fill("-e.1\r")

	_, err := control.SetPower(true)
	assert.Equal(t, errorCodes[0], err)

	mock.Fill("-e.2\r")

	_, err = control.SetPower(true)
	assert.Equal(t, errorCodes[1], err)

	mock.Fill("-e.3\r")

	_, err = control.SetPower(true)
	assert.Equal(t, errorCodes[2], err)

	mock.Fill("-e.0\r")

	_, err = control.SetPower(true)
	assert.Equal(t, errUnexpectedResponse, err)
}
